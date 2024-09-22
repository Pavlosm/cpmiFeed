package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Data        string             `bson:"data"`
	Description string             `bson:"event_type"`
	URL         string             `bson:"url"`
	Tags        []string           `bson:"tags"`
	Timestamp   time.Time          `bson:"timestamp"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
}

type UserEventFilters struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
}

type UserEventsDocument struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Events []UserEvent        `bson:"events"`
}

type UserEvent struct {
	EventID   string      `bson:"event_id"`
	EventType string      `bson:"event_type"`
	Timestamp time.Time   `bson:"timestamp"`
	Data      interface{} `bson:"data"`
}
