package repositories

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/models/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDexRepository interface {
	GetDexByFactory(factory string) (*entities.Dex, error)
	GetDexes() (*[]entities.Dex, error)
	CreateDex(dex *entities.Dex) (*mongo.InsertOneResult, error)
}

type dexRepository struct {
	dexCollection *mongo.Collection
	ctx           context.Context
}

func NewDexRepository() IDexRepository {
	return &dexRepository{
		dexCollection: db.DexCollection,
		ctx:           db.Ctx,
	}
}

func (d *dexRepository) GetDexByFactory(factory string) (*entities.Dex, error) {
	var dex entities.Dex
	err := d.dexCollection.FindOne(d.ctx, bson.M{"factory": factory}).Decode(&dex)
	if err != nil {
		return nil, err
	}
	return &dex, nil
}

func (d *dexRepository) GetDexes() (*[]entities.Dex, error) {
	var dexes []entities.Dex
	cursor, err := d.dexCollection.Find(d.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(d.ctx, &dexes)
	if err != nil {
		return nil, err
	}
	return &dexes, nil
}

func (d *dexRepository) CreateDex(dex *entities.Dex) (*mongo.InsertOneResult, error) {
	return d.dexCollection.InsertOne(d.ctx, dex)
}
