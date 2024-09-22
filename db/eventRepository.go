package db

import (
	"context"
	"cpmiFeed/rawEventModels"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository interface {
	Save(events []rawEventModels.Event) error
}

type MongoEventRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoEventRepository(uri, database string) (*MongoEventRepository, error) {

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	return &MongoEventRepository{
		client:     client,
		database:   database,
		collection: "Event",
	}, nil
}

func (r *MongoEventRepository) Save(events []rawEventModels.Event) error {

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

func (r *MongoEventRepository) ReadEvents(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]rawEventModels.Event, error) {

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

	rawEvents := make([]rawEventModels.Event, len(events))

	for _, event := range events {
		rawEvents = append(rawEvents, NewEventFromDocument(event))
	}

	return rawEvents, nil
}

func (r *MongoEventRepository) ReadEventsPaginated(ctx context.Context, filter interface{}, page, pageSize int64) ([]rawEventModels.Event, error) {
	skip := (page - 1) * pageSize

	opts := options.Find().SetSkip(skip).SetLimit(pageSize)

	return r.ReadEvents(ctx, filter, opts)
}
