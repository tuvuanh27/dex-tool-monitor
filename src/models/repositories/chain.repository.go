package repositories

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/models/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IChainRepository interface {
	CreateChain(chain *entities.Chain) error
	FindAll() ([]*entities.Chain, error)
	GetChainByName(name string) (*entities.Chain, error)
	UpdateBlock(chainName string, version, block int) *mongo.SingleResult
}

type ChainRepository struct {
	chainCollection *mongo.Collection
	ctx             context.Context
}

func NewChainRepository() IChainRepository {
	return &ChainRepository{
		chainCollection: db.ChainCollection,
		ctx:             db.Ctx,
	}
}

func (c *ChainRepository) CreateChain(chain *entities.Chain) error {
	_, err := c.chainCollection.InsertOne(c.ctx, chain)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChainRepository) FindAll() ([]*entities.Chain, error) {
	var chains []*entities.Chain
	cursor, err := c.chainCollection.Find(c.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c.ctx, &chains)
	if err != nil {
		return nil, err
	}
	return chains, nil
}

func (c *ChainRepository) GetChainByName(name string) (*entities.Chain, error) {
	var chain entities.Chain
	err := c.chainCollection.FindOne(c.ctx, bson.M{"name": name}).Decode(&chain)
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

func (c *ChainRepository) UpdateBlock(chainName string, version, block int) *mongo.SingleResult {
	return c.chainCollection.FindOneAndUpdate(
		c.ctx,
		bson.M{"name": chainName, "topic_block.version": version},
		bson.M{"$set": bson.M{"topic_block.$.block": block}},
	)
}
