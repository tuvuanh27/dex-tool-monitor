package entities

type Token struct {
	Address  string `json:"address" bson:"address"`
	Name     string `json:"name" bson:"name"`
	Symbol   string `json:"symbol" bson:"symbol"`
	Decimals int    `json:"decimals" bson:"decimals"`
}
