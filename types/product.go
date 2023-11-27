package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Rating      int                `bson:"rating" json:"rating"`
	Price       float64            `bson:"price" json:"price"`
	Description string             `bson:"description" json:"description"`
}

type UpdateProductParams struct {
	Quantity    int     `bson:"quantity" json:"quantity"`
	Rating      int     `bson:"rating" json:"rating"`
	Price       float64 `bson:"price" json:"price"`
	Description string  `bson:"description" json:"description"`
}

type CreateProductParams struct {
	Name        string  `bson:"name" json:"name"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Price       float64 `bson:"price" json:"price"`
	Description string  `bson:"description" json:"description"`
}

func CreateProductFromParams(params CreateProductParams) (*Product, error) {
	return &Product{
		Name:        params.Name,
		Quantity:    params.Quantity,
		Price:       params.Price,
		Description: params.Description,
	}, nil
}

func (c *CreateProductParams) Validate() map[string]string {
	errors := make(map[string]string)

	if len(c.Name) == 0 {
		errors["name"] = "Name is mandatory"
	}
	if len(c.Description) == 0 {
		errors["description"] = "Product description is mandatory"
	}
	if c.Price <= 0 {
		errors["price"] = fmt.Sprintf("Price must be greater than %f", c.Price)
	}
	if c.Quantity <= 0 {
		errors["quantity"] = "Quantity is mandatory"
	}
	return errors
}
