package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ItemId primitive.ObjectID `bson:"itemId" json:"itemId"`
	Count  int                `bson:"count" json:"count"`
	Price  float64            `bson:"price" json:"price"`
}

func CreateCartItemfromProduct(product Product) (*CartItem, error) {
	return &CartItem{
		ItemId: product.ID,
		Count:  1,
		Price:  product.Price,
	}, nil
}
