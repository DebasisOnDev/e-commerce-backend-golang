package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStore interface {
	GetProducts(context.Context) ([]*types.Product, error)
	GetProductByID(context.Context, string) (*types.Product, error)
	InsertProduct(context.Context, *types.Product) (*types.Product, error)
	DeleteProduct(context.Context, string) error
	UpdateProduct(context.Context, string, bson.M) error
}

type MongoProductStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoProductStore(client *mongo.Client) *MongoProductStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	return &MongoProductStore{
		client: client,
		coll:   client.Database(dbname).Collection("products"),
	}
}

func (p *MongoProductStore) GetProducts(ctx context.Context) ([]*types.Product, error) {
	cur, err := p.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var products []*types.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
func (p *MongoProductStore) GetProductByID(ctx context.Context, id string) (*types.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product types.Product
	if err := p.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (p *MongoProductStore) InsertProduct(ctx context.Context, product *types.Product) (*types.Product, error) {
	res, err := p.coll.InsertOne(ctx, product)
	if err != nil {
		fmt.Println("cant insert product")
		return nil, err
	}
	product.ID = res.InsertedID.(primitive.ObjectID)
	return product, nil
}

func (p *MongoProductStore) DeleteProduct(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	fmt.Println(oid)
	if err != nil {
		return err
	}
	res, err := p.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents", res.DeletedCount)
	return nil
}

func (p *MongoProductStore) UpdateProduct(ctx context.Context, id string, update bson.M) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	u := bson.M{"$set": update}
	_, err = p.coll.UpdateByID(ctx, oid, u)
	return err
}

// func (s *MongoCartStore) UpdateTheCart(ctx context.Context, id string, update bson.M) error {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}
// 	u := bson.M{"$set": update}
// 	_, err = s.coll.UpdateByID(ctx, oid, u)
// 	return err
// }
