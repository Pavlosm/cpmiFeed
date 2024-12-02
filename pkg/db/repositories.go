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
	UserFilter UserFilterRepository
	UserEvents UserEventsRepository
}

func (r *Repositories) Close() {
	r.Event.Close()
	r.User.Close()
	r.UserEvents.Close()
}

func NewClient() *mongo.Client {
	a := os.Getenv("MONGO_ADDRESS")
	u := os.Getenv("MONGO_USERNAME")
	p := os.Getenv("MONGO_PASSWORD")

	var connectionString string
	if u != "" && p != "" {
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s", u, p, a)
	} else {
		connectionString = fmt.Sprintf("mongodb://%s", a)
	}

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(connectionString))

	if err != nil {
		panic(err)
	}

	return client
}

func GetDatabase() string {
	return os.Getenv("MONGO_DATABASE")
}

func NewRepositories() *Repositories {

	d := GetDatabase()

	client := NewClient()

	return &Repositories{
		Event:      NewMongoEventRepository(client, d),
		User:       NewMongoUserRepository(client, d),
		UserEvents: NewMongoUserEventsRepository(client, d),
		UserFilter: NewMongoUserFilterRepository(client, d),
	}
}
