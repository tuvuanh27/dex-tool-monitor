package entities

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d *Dex) MarshalBSON() ([]byte, error) {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	d.UpdatedAt = time.Now()

	type my Dex
	return bson.Marshal((*my)(d))
}

type Dex struct {
	ID           string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Factory      string    `json:"factory" bson:"factory"`
	Router       []string  `json:"router" bson:"router"`
	FactoryBlock int64     `json:"factory_block" bson:"factory_block"`
	Name         string    `json:"name" bson:"name"`
	Chain        string    `json:"chain" bson:"chain"`
	Version      string    `json:"version" bson:"version"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
