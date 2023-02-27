package dbutils

import (
	"context"
	"fmt"
	"time"

	"github.com/Anuragch/go_auth/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDb(config *configs.DbConfig) (*mongo.Database, error) {
	fmt.Println(config)
	//uri := "mongodb://" + config.Hosts + "/" + config.DbName + "?" + config.DbOptions
	uri := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Could not establish DB connection")
		return nil, err
	}

	if err := connectMongo(client); err != nil {
		return nil, err
	}

	if err := pingMongo(client, readpref.Primary()); err != nil {
		return nil, err
	}

	return client.Database("user_db"), nil
}

func connectMongo(c *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c.Connect(ctx)
}

func pingMongo(c *mongo.Client, rp *readpref.ReadPref) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c.Ping(ctx, rp)
}
