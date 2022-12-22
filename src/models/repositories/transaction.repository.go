package repositories

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/models/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITransactionRepository interface {
	CreateTransaction(transaction *entities.Transaction) (*mongo.InsertOneResult, error)
	UpdateTransaction(transaction *entities.Transaction) (*mongo.UpdateResult, error)
}

type transactionRepository struct {
	transactionRepository *mongo.Collection
	ctx                   context.Context
}

func NewTransactionRepository() ITransactionRepository {
	return &transactionRepository{
		transactionRepository: db.TransactionCollection,
		ctx:                   db.Ctx,
	}
}

func (t *transactionRepository) CreateTransaction(transaction *entities.Transaction) (*mongo.InsertOneResult, error) {
	return t.transactionRepository.InsertOne(t.ctx, transaction)
}

func (t *transactionRepository) UpdateTransaction(transaction *entities.Transaction) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{
		{"transaction_hash", transaction.TransactionHash},
		{"transaction_index", transaction.TransactionIndex},
		{"log_index", transaction.LogIndex},
	}
	update := bson.D{
		{"$set",
			bson.D{
				{"price", transaction.Price},
				{"price_usd_0", transaction.PriceUSD0},
				{"volume_0", transaction.Volume0},
				{"volume_usd_0", transaction.VolumeUSD0},
				{"price_usd_1", transaction.PriceUSD1},
				{"volume_1", transaction.Volume1},
				{"volume_usd_1", transaction.VolumeUSD1}},
		},
	}

	return t.transactionRepository.UpdateOne(t.ctx, filter, update, opts)
}
