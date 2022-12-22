package repositories

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/models/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITokenRepository interface {
	GetTokenByAddress(address string) (*entities.Token, error)
	CreateToken(token *entities.Token) (*mongo.InsertOneResult, error)
}

type tokenRepository struct {
	tokenCollection *mongo.Collection
	ctx             context.Context
}

func NewTokenRepository() ITokenRepository {
	return &tokenRepository{
		tokenCollection: db.TokenCollection,
		ctx:             db.Ctx,
	}
}

func (t *tokenRepository) GetTokenByAddress(address string) (*entities.Token, error) {
	var token entities.Token
	err := t.tokenCollection.FindOne(t.ctx, bson.M{"address": address}).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *tokenRepository) CreateToken(token *entities.Token) (*mongo.InsertOneResult, error) {
	return t.tokenCollection.InsertOne(t.ctx, token)
}
