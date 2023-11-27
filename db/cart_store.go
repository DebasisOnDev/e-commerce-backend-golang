package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartStore interface {
	InsertToCart(context.Context, *types.CartItem) (*types.CartItem, error)
	GetFromCart(context.Context) ([]*types.CartItem, error)
	RemoveFromCart(context.Context, string) error
	EmptyTheCart(context.Context, bson.M) error
}

type MongoCartStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoCartStore(client *mongo.Client) *MongoCartStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	return &MongoCartStore{
		client: client,
		coll:   client.Database(dbname).Collection("cart"),
	}
}
func (s *MongoCartStore) EmptyTheCart(ctx context.Context, filter bson.M) error {
	_, err := s.coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoCartStore) GetFromCart(ctx context.Context) ([]*types.CartItem, error) {
	var products []*types.CartItem
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	fmt.Println(products)
	return products, nil
}

func (s *MongoCartStore) InsertToCart(ctx context.Context, item *types.CartItem) (*types.CartItem, error) {
	prod := s.coll.FindOne(ctx, bson.M{"itemId": item.ItemId})
	var existingItem types.CartItem
	if err := prod.Decode(&existingItem); err == nil {
		return nil, nil
	}

	res, err := s.coll.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}
	item.ID = res.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (s *MongoCartStore) RemoveFromCart(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s removed from cart", id)
	return nil
}
