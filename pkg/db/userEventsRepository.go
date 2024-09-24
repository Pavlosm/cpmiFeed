package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEventsRepository interface {
	CreateUserEvents(ctx context.Context, userEvents UserEvents) error
	Close() error
}

type MongoUserEventsRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoUserEventsRepository(client *mongo.Client, database string) *MongoUserEventsRepository {
	return &MongoUserEventsRepository{
		client:     client,
		database:   database,
		collection: "UserEvents"}
}

func (r *MongoUserEventsRepository) CreateUserEvents(ctx context.Context, userEvents UserEvents) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.InsertOne(ctx, userEvents)
	return err
}

func (r *MongoUserEventsRepository) CreateUserEventsWithID(ctx context.Context, id primitive.ObjectID, userEvents UserEvents) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	userEvents.ID = id
	_, err := coll.InsertOne(ctx, userEvents)
	return err
}

func (r *MongoUserEventsRepository) GetUserEvents(ctx context.Context, userID primitive.ObjectID) (UserEvents, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": userID}
	var userEvents UserEvents
	err := coll.FindOne(ctx, filter).Decode(&userEvents)
	return userEvents, err
}

func (r *MongoUserEventsRepository) UpdateUserEvents(ctx context.Context, userEvents UserEvents) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": userEvents.ID}, bson.M{"$set": userEvents})
	return err
}

func (r *MongoUserEventsRepository) UpdateSingleEvent(ctx context.Context, userID primitive.ObjectID, event UserEvent) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": userID, "events.event_id": event.EventID}
	update := bson.M{"$set": bson.M{
		"events.$.event_type": event.EventType,
		"events.$.timestamp":  event.Timestamp,
		"events.$.viewed":     event.Viewed,
		"events.$.deleted":    event.Deleted,
	}}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUserEventsRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}
