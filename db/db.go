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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
)

type Connection struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
	Context    context.Context
}

func Connect(config *configuration.MongoConfig) *Connection {
	uri := getUri(config.Username, config.Password, config.Hostname, config.RetryWrites)
	database := config.Database
	collection := config.Collection
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	utils.Check(err)

	db := client.Database(database)

	return &Connection{
		Client:     client,
		Database:   db,
		Collection: db.Collection(collection),
		Context:    context.TODO(),
	}
}

func (c *Connection) Close() {
	if err := c.Client.Disconnect(c.Context); err != nil {
		panic(err)
	}
}

func (c *Connection) findOne(id string) error {
	var result interface{}
	filter := bson.D{{"_id", id}}
	return c.Collection.FindOne(context.TODO(), filter, nil).Decode(&result)
}

func (c *Connection) InsertBatch(documents []interface{}) {
	if len(documents) > 0 {
		result, err := c.Collection.InsertMany(context.TODO(), documents, nil)
		utils.Check(err)
		if result != nil {
			fmt.Printf("New batch with %d\n", len(documents))
		}
	}
}

func (c *Connection) GetConditionalDocument(data *domain.JobPost) bson.D {
	id := generateHashID(data)
	err := c.findOne(id)
	if err != nil {
		metadata := bson.D{
			{"company", data.Company},
			{"position", data.Position},
			{"location", data.Location},
			{"benefit", data.Benefit},
			{"Link", data.Link},
		}
		document := bson.D{
			{"_id", id},
			{"timestamp", data.PostDate},
			{"metadata", metadata},
			{"jobs_count", 1},
		}
		log.Printf("New Job [%s] for [%s] %s\n", data.Position, data.Company, data.RawPostDate)
		return document
	}
	log.Printf("Found in DB! - Job [%s] for [%s] %s\n", data.Position, data.Company, data.RawPostDate)
	return nil
}

func generateHashID(jobPost *domain.JobPost) string {
	seed := fmt.Sprintf("%s_%s_%s_%s", jobPost.Position, jobPost.Company, jobPost.Location, jobPost.RawPostDate)
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
