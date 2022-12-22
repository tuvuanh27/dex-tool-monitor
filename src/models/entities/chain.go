package entities

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *Chain) MarshalBSON() ([]byte, error) {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.UpdatedAt = time.Now()

	type my Chain
	return bson.Marshal((*my)(c))
}

type TopicBlock struct {
	Topic   string `json:"topic" bson:"topic"`
	Version int    `json:"version" bson:"version"`
	Block   int  `json:"block" bson:"block"`
}

type Chain struct {
	ID         string       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string       `json:"name" bson:"name"`
	HttpRpc    string       `json:"http_rpc" bson:"http_rpc"`
	WsRpc      string       `json:"ws_rpc" bson:"ws_rpc"`
	TopicBlock []TopicBlock `json:"topic_block" bson:"topic_block"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" bson:"updated_at"`
}
