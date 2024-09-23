package db

import (
	"context"
	"cpmiFeed/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, user User) (common.User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, user User) error
	Close() error
}

type MongoUserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoUserRepository(client *mongo.Client, database string) *MongoUserRepository {
	return &MongoUserRepository{
		client:     client,
		database:   database,
		collection: "User",
	}
}

func (r *MongoUserRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}

func (r *MongoUserRepository) GetAll(ctx context.Context) ([]User, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoUserRepository) CreateUser(ctx context.Context, user User) (common.User, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	v, err := coll.InsertOne(ctx, user)
	if err != nil {
		return common.User{}, err
	}
	user.ID = v.InsertedID.(primitive.ObjectID)
	return NewUserFromDocument(user), nil
}

func (r *MongoUserRepository) GetById(ctx context.Context, id string) (User, error) {
	coll := r.client.Database(r.database).Collection(r.collection)
	var user User
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *MongoUserRepository) UpdateUser(ctx context.Context, user User) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (r *MongoUserRepository) DeleteUser(ctx context.Context, user User) error {
	coll := r.client.Database(r.database).Collection(r.collection)
	_, err := coll.DeleteOne(ctx, bson.M{"_id": user.ID})
	return err
}
