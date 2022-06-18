package db

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/margostino/job-pulse/configuration"
	"github.com/margostino/job-pulse/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"time"
)

func BuildJobPost(jobTextParts []string, link string, rawPostDate string, postDate time.Time) *JobPost {
	return &JobPost{
		Position:    utils.GetOrDefault(0, jobTextParts),
		Company:     utils.GetOrDefault(1, jobTextParts),
		Location:    utils.GetOrDefault(2, jobTextParts),
		Benefit:     utils.GetOrDefault(3, jobTextParts),
		Link:        link,
		RawPostDate: rawPostDate,
		PostDate:    postDate,
	}
}

func generateHashID(jobPost *JobPost) string {
	seed := fmt.Sprintf("%s_%s_%s_%s", jobPost.Position, jobPost.Company, jobPost.Location, jobPost.RawPostDate)
	return HashFrom(seed)
}

func findOne(collection *mongo.Collection, id string) error {
	var result interface{}
	filter := bson.D{{"_id", id}}
	return collection.FindOne(context.TODO(), filter, nil).Decode(&result)
}

func InsertBatch(collection *mongo.Collection, documents []interface{}) {
	result, err := collection.InsertMany(context.TODO(), documents, nil)
	utils.Check(err)
	if result != nil {
		fmt.Printf("New batch with %d\n", len(documents))
	}
}

func ConditionalInsert(collection *mongo.Collection, jobPost *JobPost) bson.D {
	id := generateHashID(jobPost)
	err := findOne(collection, id)
	if err != nil {
		metadata := bson.D{
			{"company", jobPost.Company},
			{"position", jobPost.Position},
			{"location", jobPost.Location},
			{"benefit", jobPost.Benefit},
			{"Link", jobPost.Link},
		}
		document := bson.D{
			{"_id", id},
			{"timestamp", jobPost.PostDate},
			{"metadata", metadata},
			{"jobs_count", 1},
		}
		return document
		log.Printf("New Job [%s] for [%s] %s\n", jobPost.Position, jobPost.Company, jobPost.RawPostDate)
	}
	return nil
}

func ConnectCollection(config *configuration.MongoConfig) *mongo.Collection {
	uri := getUri(config.Username, config.Password, config.Hostname, config.RetryWrites)
	database := config.Database
	collection := config.Collection
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	return client.Database(database).Collection(collection)
}

func getUri(username string, password string, hostname string, retryWrites bool) string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=%t&w=majority", username, password, hostname, retryWrites)
}

func HashFrom(seed string) string {
	hash := sha1.New()
	io.WriteString(hash, seed)
	return hex.EncodeToString(hash.Sum(nil))
}
