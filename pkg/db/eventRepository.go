package db

import (
	"context"
	"cpmiFeed/pkg/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	Save(events []common.Event) error
	ReadEvents(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]common.Event, error)
	GetById(ctx context.Context, id string) (common.Event, error)
	Close() error
}

type MongoEventRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoEventRepository(client *mongo.Client, database string) *MongoEventRepository {

	return &MongoEventRepository{
		client:     client,
		database:   database,
		collection: EventsCollectionName,
	}
}

func (r *MongoEventRepository) Save(events []common.Event) error {

	coll := r.client.Database(r.database).Collection(r.collection)

	docs := make([]interface{}, len(events))

	for i, event := range events {
		docs[i] = NewEventDocument(event)
	}

	_, err := coll.InsertMany(context.TODO(), docs)

	return err
}

func (r *MongoEventRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}

func (r *MongoEventRepository) ReadEvents(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]common.Event, error) {

	coll := r.client.Database(r.database).Collection(r.collection)

	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []Event
	if err = cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	var rawEvents []common.Event

	for _, event := range events {
		rawEvents = append(rawEvents, NewEventFromDocument(event))
	}

	return rawEvents, nil
}

func (r *MongoEventRepository) GetById(ctx context.Context, id string) (common.Event, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	filter := bson.M{"_id": id}
	var event Event
	err := coll.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return common.Event{}, err
	}
	return NewEventFromDocument(event), nil
}

func (r *MongoEventRepository) ReadEventsPaginated(ctx context.Context, filter interface{}, page, pageSize int64) ([]common.Event, error) {
	skip := (page - 1) * pageSize

	opts := options.Find().SetSkip(skip).SetLimit(pageSize)

	return r.ReadEvents(ctx, filter, opts)
}
