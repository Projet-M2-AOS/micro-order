package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/b4cktr4ck5r3/micro-order/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var MI MongoInstance

func ConnectDB() {
	credential := options.Credential{
		Username: config.Config("MONGO_USER"),
		Password: config.Config("MONGO_PASSWORD"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.Config("MONGO_URI")).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected!")

	MI = MongoInstance{
		Client: client,
		DB:     client.Database(config.Config("DB")),
	}
}
