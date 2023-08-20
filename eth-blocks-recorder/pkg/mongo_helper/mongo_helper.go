package mongo_helper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectToMongo(ethBlocksMongo string, mongoUser string, ethBlocksMongoPassword string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(ethBlocksMongo)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUser,
		Password: ethBlocksMongoPassword,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("eth-blocks-recorder::Error connecting:", err)
		return nil, err
	}

	return c, nil
}
