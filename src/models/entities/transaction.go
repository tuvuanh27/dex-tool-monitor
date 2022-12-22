package entities

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TransactionType int

const (
	Buy TransactionType = iota
	Sell
)

func (t *Transaction) MarshalBSON() ([]byte, error) {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	t.UpdatedAt = time.Now()

	type my Transaction
	return bson.Marshal((*my)(t))
}

type Transaction struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PairAddress      string             `json:"pair_address" bson:"pair_address"`
	DexId            string             `json:"dex_id" bson:"dex_id"`
	Type             TransactionType    `json:"type" bson:"type"`
	Amount0In        string             `json:"amount0_in" bson:"amount0_in"`
	Amount1In        string             `json:"amount1_in" bson:"amount1_in"`
	Amount0Out       string             `json:"amount0_out" bson:"amount0_out"`
	Amount1Out       string             `json:"amount1_out" bson:"amount1_out"`
	Amount0          string             `json:"amount0" bson:"amount0"`
	Amount1          string             `json:"amount1" bson:"amount1"`
	BlockNumber      int64              `json:"block_number" bson:"block_number"`
	Timestamp        int64              `json:"timestamp" bson:"timestamp"`
	TransactionHash  string             `json:"transaction_hash" bson:"transaction_hash"`
	TransactionIndex int                `json:"transaction_index" bson:"transaction_index"`
	LogIndex         int                `json:"log_index" bson:"log_index"`
	Price            float64            `json:"price" bson:"price"`
	PriceUSD0        float64            `json:"price_usd_0" bson:"price_usd_0"`
	Volume0          float64            `json:"volume_0" bson:"volume_0"`
	VolumeUSD0       float64            `json:"volume_usd_0" bson:"volume_usd_0"`
	PriceUSD1        float64            `json:"price_usd_1" bson:"price_usd_1"`
	Volume1          float64            `json:"volume_1" bson:"volume_1"`
	VolumeUSD1       float64            `json:"volume_usd_1" bson:"volume_usd_1"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}
