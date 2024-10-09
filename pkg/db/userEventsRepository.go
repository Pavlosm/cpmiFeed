package db

import (
	"context"
	"cpmiFeed/pkg/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserEventsRepository interface {
	UpsertUserEvents(ctx context.Context, userId string, userEvents ...common.Event) error
	GetUserEvents(ctx context.Context, userId string) ([]common.UserEvent, error)
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

func (r *MongoUserEventsRepository) UpsertUserEvents(ctx context.Context, userId string, userEvents ...common.Event) error {
	coll := r.client.Database(r.database).Collection(r.collection)

	ue, err := NewDocumentFromUserEvents(userId, userEvents)
	if err != nil {
		return err
	}

	_, err = coll.InsertOne(ctx, ue)
	if err == nil {
		return nil
	}

	v, ok := err.(mongo.WriteException)
	if !ok {
		return err
	}
	b := true
	for _, e := range v.WriteErrors {
		if e.Code == 11000 {
			b = true
			break
		}
	}

	if !b {
		return err
	}

	f := bson.M{"_id": ue.UserID}
	u := bson.M{"$push": bson.M{"events": bson.M{"$each": ue.Events}}}
	o := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, f, u, o)

	return err
}

func (r *MongoUserEventsRepository) CreateUserEventsWithID(ctx context.Context, id primitive.ObjectID, userEvents UserEvents) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	userEvents.UserID = id
	_, err := coll.InsertOne(ctx, userEvents)
	return err
}

func (r *MongoUserEventsRepository) GetUserEvents(ctx context.Context, userId string) ([]common.UserEvent, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return make([]common.UserEvent, 0), err
	}
	filter := bson.M{"_id": id}
	var userEvents UserEvents
	err = coll.FindOne(ctx, filter).Decode(&userEvents)
	if err != nil {
		return make([]common.UserEvent, 0), err
	}

	ev := NewUserEventsFromUserEventsDocument(userEvents)

	return ev, err
}

func (r *MongoUserEventsRepository) UpdateUserEvents(ctx context.Context, userEvents UserEvents) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": userEvents.UserID}, bson.M{"$set": userEvents})
	return err
}

func (r *MongoUserEventsRepository) UpdateSingleEvent(ctx context.Context, userID primitive.ObjectID, event UserEvent) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": userID, "events.event_id": event.EventID}
	update := bson.M{"$set": bson.M{
		"events.$.timestamp": event.Timestamp,
		"events.$.viewed":    event.Viewed,
		"events.$.deleted":   event.Deleted,
	}}
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoUserEventsRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}
