package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	UserID               primitive.ObjectID   `bson:"userId" json:"userId"`
	Products             []primitive.ObjectID `bson:"products" json:"products"`
	TotalPrice           float64              `bson:"totalPrice" json:"totalPrice"`
	OrderDate            time.Time            `bson:"orderDate" json:"orderDate"`
	ExpectedDeliveryDate time.Time            `bson:"expectedDate" json:"expectedDate"`
	Cancelled            bool                 `bson:"cancelled" json:"cancelled"`
}
