package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gql-with-go/graph/model"
	"log"
	"os"
	"time"
)

type VideoRepository interface {
	Save(video *model.Video)
	FindAll() []*model.Video
}

type database struct {
	client *mongo.Client
}

const (
	DATABASE   = "graphql"
	COLLECTION = "videos"
)

func New() VideoRepository {

	// mongo+srv://USERNAME:PASSWORD@HOST:PORT
	MONGODB := os.Getenv("MONGODB")

	clientOption := options.Client().ApplyURI(MONGODB)

	clientOption = clientOption.SetMaxPoolSize(50)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	dbClient, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to MongoDB!")

	return &database{
		client: dbClient,
	}
}

func (db *database) Save(video *model.Video) {
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	_, err := collection.InsertOne(context.TODO(), video)
	if err != nil {
		log.Fatalln(err)
	}

}

func (db *database) FindAll() []*model.Video {
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln(err)
	}
	defer cursor.Close(context.TODO())

	var result []*model.Video
	for cursor.Next(context.TODO()) {
		var v *model.Video
		err := cursor.Decode(&v)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, v)
	}

	return result
}
