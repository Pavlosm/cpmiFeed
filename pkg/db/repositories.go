package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	Event      EventRepository
	User       UserRepository
	UserEvents UserEventsRepository
}

func (r *Repositories) Close() {
	r.Event.Close()
	r.User.Close()
	r.UserEvents.Close()
}

func NewRepositories() *Repositories {

	a := os.Getenv("MONGO_ADDRESS")
	u := os.Getenv("MONGO_USERNAME")
	p := os.Getenv("MONGO_PASSWORD")
	d := os.Getenv("MONGO_DATABASE")

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s", u, p, a)

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(connectionString))

	if err != nil {
		panic(err)
	}

	return &Repositories{
		Event:      NewMongoEventRepository(client, d),
		User:       NewMongoUserRepository(client, d),
		UserEvents: NewMongoUserEventsRepository(client, d),
	}
}
