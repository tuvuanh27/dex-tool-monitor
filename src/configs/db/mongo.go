package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Ctx                   = context.TODO()
	PairCollection        *mongo.Collection
	TransactionCollection *mongo.Collection
	TokenCollection       *mongo.Collection
	DexCollection         *mongo.Collection
	ChainCollection       *mongo.Collection
)

func Setup(uri string, dbName string) {
	option := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(Ctx, option)
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(Ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	db := client.Database(dbName)
	PairCollection = db.Collection("pairs")
	TransactionCollection = db.Collection("transactions")
	TokenCollection = db.Collection("tokens")
	DexCollection = db.Collection("dexes")
	ChainCollection = db.Collection("chains")

	// create index for pair collection
	indexPair := []mongo.IndexModel{
		{
			Keys:    bson.D{{"address", 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err = PairCollection.Indexes().CreateMany(Ctx, indexPair)
	if err != nil {
		panic(err)
	}

	// create index for token collection
	indexToken := []mongo.IndexModel{
		{
			Keys:    bson.D{{"address", 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err = TokenCollection.Indexes().CreateMany(Ctx, indexToken)
	if err != nil {
		panic(err)
	}

	// create index for transaction collection
	indexTransaction := []mongo.IndexModel{
		{
			Keys: bson.D{{"pair_address", 1}},
		},
		{
			Keys:    bson.D{{"transaction_hash", 1}, {"transaction_index", 1}, {"log_index", 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err = TransactionCollection.Indexes().CreateMany(Ctx, indexTransaction)
	if err != nil {
		panic(err)
	}
}
