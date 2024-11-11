package db

import (
	"context"
	"cpmiFeed/pkg/common"
	"log"

	"github.com/gocql/gocql"
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

type CassandraEventRepository struct {
	cluster *gocql.ClusterConfig
}

func (r *CassandraEventRepository) Save(events []common.Event) error {
	session, err := r.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	for _, event := range events {
		if err := session.Query(`INSERT INTO events (id, data, url, description, tags, timestamp) VALUES (?, ?, ?, ?, ?, ?)`, event.ID, event.Description, event.URL, event.Description, event.Tags, event.Timestamp).Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *CassandraEventRepository) Close() error {
	return nil
}

func (r *CassandraEventRepository) ReadEvents(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]common.Event, error) {

	session, err := r.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Assuming filter is a map with the query parameters
	query := `SELECT * FROM events`
	// Execute the query with the filter parameters
	iter := session.Query(query, filter).Iter()
	defer iter.Close()

	var id string
	var data string
	var url string
	var description string
	var tags []string

	//Timestamp   time.Time `json:"timestamp"`
	events := make([]common.Event, 0)
	for iter.Scan(&id, &data, &url, &description, &tags) {
		events = append(events, common.Event{
			ID:          id,
			Data:        data,
			Description: description,
			Tags:        tags,
			URL:         url,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *CassandraEventRepository) GetById(ctx context.Context, id string) (common.Event, error) {
	session, err := r.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	return common.Event{}, nil
}

func (r *CassandraEventRepository) ReadEventsPaginated(ctx context.Context, filter interface{}, page, pageSize int64) ([]common.Event, error) {
	session, err := r.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	// skip := (page - 1) * pageSize

	// opts := options.Find().SetSkip(skip).SetLimit(pageSize)

	// return r.ReadEvents(ctx, filter, opts)
	return nil, nil
}

func NewCassandraRepository(cluster *gocql.ClusterConfig) *CassandraEventRepository {
	return &CassandraEventRepository{
		cluster: cluster,
	}
}
