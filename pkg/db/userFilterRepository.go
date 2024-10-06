package db

import (
	"context"
	"cpmiFeed/pkg/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserFilterRepository interface {
	GetAll(ctx context.Context) ([]UserEventFilters, error)
	GetForUser(ctx context.Context, userId string) ([]common.UserEventFilter, error)
	Create(ctx context.Context, userId string, filter []common.UserEventFilter) error
	Update(ctx context.Context, userId string, filter []common.UserEventFilter) error
	Delete(ctx context.Context, userId string) error
}

type MongoUserFilterRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoUserFilterRepository(client *mongo.Client, database string) UserFilterRepository {
	return &MongoUserFilterRepository{
		client:     client,
		database:   database,
		collection: "UserEventFilters",
	}
}

func (r *MongoUserFilterRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}

func (r *MongoUserFilterRepository) GetAll(ctx context.Context) ([]UserEventFilters, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var filters []UserEventFilters
	if err = cursor.All(ctx, &filters); err != nil {
		return nil, err
	}

	return filters, nil
}

func (r *MongoUserFilterRepository) GetForUser(ctx context.Context, userId string) ([]common.UserEventFilter, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	var filter UserEventFilters
	if err := coll.FindOne(ctx, bson.M{"user_id": userId}).Decode(&filter); err != nil {
		return make([]common.UserEventFilter, 0), err
	}
	f := NewFiltersFromDocument(filter)
	return f, nil
}

func (r *MongoUserFilterRepository) Create(ctx context.Context, userId string, filters []common.UserEventFilter) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	doc, err := NewDocumentFromFilters(userId, filters)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, doc)
	return err
}

func (r *MongoUserFilterRepository) Update(ctx context.Context, userId string, filters []common.UserEventFilter) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	doc := NewDocFiltersFromCommonFilters(filters)
	_, err := coll.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": bson.M{"filters": doc}})
	return err
}

func (r *MongoUserFilterRepository) Delete(ctx context.Context, userId string) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.DeleteOne(ctx, bson.M{"user_id": userId})
	return err
}
