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

//Mongo instance
type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var MI MongoInstance

//Make connection to Mongo database
func ConnectDB() {
	//Credentials retrieved from env variables
	credential := options.Credential{
		Username: config.Config("MONGO_USER"),
		Password: config.Config("MONGO_PASSWORD"),
	}

	//Make client
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Config("MONGO_URI")).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}

	//Set database context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Connect client to database
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Check if connection is established
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	//Output connection if success
	fmt.Println("Database connected!")

	//Set mongo instance to current connected client and database
	MI = MongoInstance{
		Client: client,
		DB:     client.Database(config.Config("DB")),
	}
}
