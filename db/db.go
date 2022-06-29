package db

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/domain"
	"github.com/margostino/job-pulse/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"time"
)

type Connection struct {
	Client             *mongo.Client
	Database           *mongo.Database
	JobPostsCollection *mongo.Collection
	GeoCollection      *mongo.Collection
	BatchesCollection  *mongo.Collection
	Context            context.Context
}

func Connect(config *configuration.Configuration) *Connection {
	uri := getUri(config.Mongo.Username, config.Mongo.Password, config.Mongo.Hostname, config.Mongo.RetryWrites)
	database := config.Mongo.Database
	jobPostsCollection := config.Mongo.JobPostsCollection
	geoCollection := config.Mongo.GeocodingCollection
	batchesCollection := config.Mongo.BatchesCollection
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	utils.Check(err)

	db := client.Database(database)

	return &Connection{
		Client:             client,
		Database:           db,
		JobPostsCollection: db.Collection(jobPostsCollection),
		GeoCollection:      db.Collection(geoCollection),
		BatchesCollection:  db.Collection(batchesCollection),
		Context:            context.TODO(),
	}
}

func (c *Connection) Close() {
	if err := c.Client.Disconnect(c.Context); err != nil {
		panic(err)
	}
}

func (c *Connection) findOneJobBy(id string) (interface{}, error) {
	var result interface{}
	filter := bson.D{{"_id", id}}
	err := c.JobPostsCollection.FindOne(context.TODO(), filter, nil).Decode(&result)
	return result, err
}

func (c *Connection) findOneGeoBy(id string) (primitive.M, error) {
	var geocoding Geocoding
	filter := bson.D{{"_id", id}}
	err := c.GeoCollection.FindOne(context.TODO(), filter, nil).Decode(&geocoding)
	if err == nil {
		geocodingMap := geocoding.Geocoding.(primitive.D).Map()
		return geocodingMap, nil
	}
	return nil, err
}

func (c *Connection) InsertOneGeocoding(location string, geocoding interface{}) {
	id := hashFrom(location)
	document := bson.D{
		{"_id", id},
		{"geocoding", geocoding},
	}
	_, err := c.GeoCollection.InsertOne(context.TODO(), document, nil)
	utils.Check(err)
}

func (c *Connection) InsertBatch(documents []interface{}, stats *domain.Stats) {
	var status string
	if len(documents) > 0 {
		result, err := c.JobPostsCollection.InsertMany(context.TODO(), documents, nil)
		if err != nil {
			status = "failed"
		}
		utils.Check(err)
		if result != nil {
			fmt.Printf("New batch with %d\n", len(documents))
			status = "ok"
		} else {
			status = "no-result"
		}
	}
	c.InsertBatchStats(status, len(documents), stats)
}

func (c *Connection) InsertBatchStats(status string, total int, stats *domain.Stats) {
	duration := time.Now().UTC().Sub(stats.StartTime)
	metadata := bson.D{
		{"jobs_count", total},
		{"status", status},
		{"duration", duration.String()},
		{"position_input", stats.PositionInput},
		{"location_input", stats.LocationInput},
	}
	document := bson.D{
		{"timestamp", time.Now().UTC()},
		{"metadata", metadata},
	}
	_, err := c.BatchesCollection.InsertOne(context.TODO(), document, nil)
	utils.Check(err)
}

func (c *Connection) FindOneGeoBy(location string) (primitive.M, error) {
	id := hashFrom(location)
	return c.findOneGeoBy(id)
}

func (c *Connection) GetGeoDocument(geocoding interface{}) bson.D {
	return bson.D{
		{"geocoding", geocoding},
	}
}

func (c *Connection) GetConditionalDocument(data *domain.JobPost) bson.D {
	id := generateHashID(data)
	_, err := c.findOneJobBy(id)
	if err != nil {
		coordinates := bson.A{data.Longitude, data.Latitude}
		geo := bson.D{
			{"type", "Point"},
			{"coordinates", coordinates},
		}
		metadata := bson.D{
			{"company", data.Company},
			{"position", data.Position},
			{"location", data.Location},
			{"geo", geo},
			{"benefit", data.Benefit},
			{"Link", data.Link},
		}
		document := bson.D{
			{"_id", id},
			{"timestamp", data.PostDate},
			{"metadata", metadata},
			{"jobs_count", 1},
		}
		log.Printf("New Job [%s] for [%s - %s] %s\n", data.Position, data.Company, data.Location, data.RawPostDate)
		return document
	}
	//log.Printf("Found in DB! - Job [%s] for [%s - %s] %s\n", data.Position, data.Company, data.Location, data.RawPostDate)
	return nil
}

func generateHashID(jobPost *domain.JobPost) string {
	seed := fmt.Sprintf("%s_%s_%s_%s", jobPost.Position, jobPost.Company, jobPost.Location, jobPost.Link)
	return hashFrom(seed)
}

func getUri(username string, password string, hostname string, retryWrites bool) string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=%t&w=majority", username, password, hostname, retryWrites)
}

func hashFrom(seed string) string {
	hash := sha1.New()
	io.WriteString(hash, seed)
	return hex.EncodeToString(hash.Sum(nil))
}
