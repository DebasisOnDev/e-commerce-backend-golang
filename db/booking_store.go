package db

import (
	"context"
	"errors"
	"os"

	"github.com/DebasisOnDev/E-Commerce-Backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	CartStore
}

func NewMongoBookingStore(client *mongo.Client, cartStore CartStore) *MongoBookingStore {
	dbname := os.Getenv("MONGO_DB_NAME")
	return &MongoBookingStore{
		client:    client,
		coll:      client.Database(dbname).Collection("bookings"),
		CartStore: cartStore,
	}
}

func (b *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := b.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (b *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := b.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &booking, nil
}

func (b *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	bookingitem, err := b.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = bookingitem.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (b *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	u := bson.M{"$set": update}
	_, err = b.coll.UpdateByID(ctx, oid, u)
	return err
}
