package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	Event      EventRepository
	User       UserRepository
	UserEvents UserEventsRepository
}

func NewRepositories(connectionString, dbName string) *Repositories {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(connectionString))

	if err != nil {
		panic(err)
	}

	return &Repositories{
		Event:      NewMongoEventRepository(client, dbName),
		User:       NewMongoUserRepository(client, dbName),
		UserEvents: NewMongoUserEventsRepository(client, dbName),
	}
}

func (r *Repositories) Close() {
	r.Event.Close()
	r.User.Close()
	r.UserEvents.Close()
}
