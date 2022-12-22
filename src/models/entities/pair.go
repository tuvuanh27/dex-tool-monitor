package entities

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (p *Pair) MarshalBSON() ([]byte, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()

	type my Pair
	return bson.Marshal((*my)(p))
}

type Pair struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Address        string             `json:"address" bson:"address"`
	DexId          string             `json:"dex_id" bson:"dex_id"`
	Token0Address  string             `json:"token_0_address" bson:"token_0_address"`
	Token1Address  string             `json:"token_1_address" bson:"token_1_address"`
	Token0Symbol   string             `json:"token_0_symbol" bson:"token_0_symbol"`
	Token1Symbol   string             `json:"token_1_symbol" bson:"token_1_symbol"`
	Token0Decimals int                `json:"token_0_decimals" bson:"token_0_decimals"`
	Token1Decimals int                `json:"token_1_decimals" bson:"token_1_decimals"`
	Price          float64            `json:"price" bson:"price"`
	PriceUsd0      float64            `json:"price_usd_0" bson:"price_usd_0"`
	PriceUsd1      float64            `json:"price_usd_1" bson:"price_usd_1"`

	Txns24H int `json:"txns_24h" bson:"txns_24h"`
	Txns6H  int `json:"txns_6h" bson:"txns_6h"`
	Txns1H  int `json:"txns_1h" bson:"txns_1h"`
	Txns5M  int `json:"txns_5m" bson:"txns_5m"`

	Buy24H float64 `json:"buy_24h" bson:"buy_24h"`
	Buy6H  float64 `json:"buy_6h" bson:"buy_6h"`
	Buy1H  float64 `json:"buy_1h" bson:"buy_1h"`
	Buy5M  float64 `json:"buy_5m" bson:"buy_5m"`

	Sell24H float64 `json:"sell_24h" bson:"sell_24h"`
	Sell6H  float64 `json:"sell_6h" bson:"sell_6h"`
	Sell1H  float64 `json:"sell_1h" bson:"sell_1h"`
	Sell5M  float64 `json:"sell_5m" bson:"sell_5m"`

	Volume24H float64 `json:"volume_24h" bson:"volume_24h"`
	Volume6H  float64 `json:"volume_6h" bson:"volume_6h"`
	Volume1H  float64 `json:"volume_1h" bson:"volume_1h"`
	Volume5M  float64 `json:"volume_5m" bson:"volume_5m"`

	PriceChange24H float64 `json:"price_change_24h" bson:"price_change_24h"`
	PriceChange6H  float64 `json:"price_change_6h" bson:"price_change_6h"`
	PriceChange1H  float64 `json:"price_change_1h" bson:"price_change_1h"`
	PriceChange5M  float64 `json:"price_change_5m" bson:"price_change_5m"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
