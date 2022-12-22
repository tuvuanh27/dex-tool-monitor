package repositories

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/models/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IPairRepository interface {
	CreatePair(pair *entities.Pair) (*mongo.InsertOneResult, error)
	GetPairByAddress(address string) (*entities.Pair, error)
	GetPairs() ([]entities.Pair, error)
	GetPriceToken(tokenAddress, pairAddress string) (float64, error)
	UpdatePrice(price, priceUsd0, priceUsd1 float64, address string) *mongo.SingleResult
}

type pairRepository struct {
	pairCollection *mongo.Collection
	ctx            context.Context
}

func NewPairRepository() IPairRepository {
	return &pairRepository{
		pairCollection: db.PairCollection,
		ctx:            db.Ctx,
	}
}

func (p *pairRepository) CreatePair(pair *entities.Pair) (*mongo.InsertOneResult, error) {
	return p.pairCollection.InsertOne(p.ctx, pair)
}

func (p *pairRepository) GetPairByAddress(address string) (*entities.Pair, error) {
	var pair entities.Pair
	err := p.pairCollection.FindOne(p.ctx, bson.M{"address": address}).Decode(&pair)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (p *pairRepository) GetPairs() ([]entities.Pair, error) {
	var pairs []entities.Pair

	cur, err := p.pairCollection.Find(p.ctx, bson.D{})
	if err != nil {
		defer func(cur *mongo.Cursor, ctx context.Context) {
			err := cur.Close(ctx)
			if err != nil {

			}
		}(cur, db.Ctx)
		return pairs, err
	}

	for cur.Next(db.Ctx) {
		var pair entities.Pair
		err := cur.Decode(&pair)
		if err != nil {
			panic(err)
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (p *pairRepository) GetPriceToken(tokenAddress, pairAddress string) (float64, error) {
	var pair entities.Pair
	opts := options.FindOne().SetSort(bson.D{{"updated_at", 1}})

	err := p.pairCollection.FindOne(p.ctx, bson.D{
		{"token_0_address", tokenAddress},
		{"price_usd_0", bson.D{{"$ne", 0}}},
		{"address", bson.D{{"$ne", pairAddress}}},
	}, opts).Decode(&pair)
	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	if pair.PriceUsd0 > 0 {
		return pair.PriceUsd0, nil
	}

	err = p.pairCollection.FindOne(p.ctx, bson.D{
		{"token_1_address", tokenAddress},
		{"price_usd_1", bson.D{{"$ne", 0}}},
		{"address", bson.D{{"$ne", pairAddress}}},
	}, opts).Decode(&pair)

	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}
	if pair.PriceUsd1 > 0 {
		return pair.PriceUsd1, nil
	}
	return 0, nil
}

func (p *pairRepository) UpdatePrice(price, priceUsd0, priceUsd1 float64, address string) *mongo.SingleResult {
	return p.pairCollection.FindOneAndUpdate(p.ctx, bson.M{"address": address}, bson.M{"$set": bson.M{
		"price":       price,
		"price_usd_0": priceUsd0,
		"price_usd_1": priceUsd1,
		"updated_at":  time.Now(),
	}})
}
