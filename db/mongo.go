package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mayur-tolexo/mysqlVsmongo/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	ctx := context.TODO()
	config := common.GetConfig()

	mongoConfig := config.Database.Mongo

	if mongoConfig.Hostname == "" {
		mongoConfig.Hostname = "localhost"
	}
	if mongoConfig.Port == "" {
		mongoConfig.Port = "27017"
	}
	if mongoConfig.Database == "" {
		mongoConfig.Database = "test"
	}
	if mongoConfig.Collection == "" {
		mongoConfig.Collection = "test"
	}

	uri := fmt.Sprintf("mongodb://%v:%v", mongoConfig.Hostname, mongoConfig.Port)
	// fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(mongoConfig.Database).Collection(mongoConfig.Collection)
}

// GetMongoCollection will return the mongo collection
func GetMongoCollection() *mongo.Collection {
	return collection
}
