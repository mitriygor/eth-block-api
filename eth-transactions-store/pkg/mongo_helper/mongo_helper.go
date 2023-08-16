package mongo_helper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectToMongo(mongoUrl string, mongoUser string, mongoPassword string) (*mongo.Client, error) {

	log.Printf("EthTransactionsStore::ConnectToMongo::mongoUrl: %v\n", mongoUrl)
	log.Printf("EthTransactionsStore::ConnectToMongo::mongoUser: %v\n", mongoUser)
	log.Printf("EthTransactionsStore::ConnectToMongo::mongoPassword: %v\n", mongoPassword)

	clientOptions := options.Client().ApplyURI(mongoUrl)
	log.Printf("EthTransactionsStore::ConnectToMongo::clientOptions: %v\n", clientOptions)

	clientOptions.SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPassword,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)

	log.Printf("EthTransactionsStore::ConnectToMongo::c: %v\n", c)
	log.Printf("EthTransactionsStore::ConnectToMongo::err: %v\n", err)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return c, nil
}
