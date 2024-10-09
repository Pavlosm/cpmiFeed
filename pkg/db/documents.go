package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID          string    `bson:"_id,omitempty"`
	Data        string    `bson:"data"`
	Description string    `bson:"event_type"`
	URL         string    `bson:"url"`
	Tags        []string  `bson:"tags"`
	Timestamp   time.Time `bson:"timestamp"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
}

type UserEventFilters struct {
	UserID  primitive.ObjectID `bson:"_id,omitempty"`
	Filters []UserEventFilter  `bson:"filters"`
}

type UserEventFilter struct {
	Name   string   `bson:"name"`
	Tags   []string `bson:"tags"`
	Tokens []string `bson:"tokens"`
}

type UserEvents struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Events []UserEvent        `bson:"events"`
}

type UserEvent struct {
	EventID   string    `bson:"event_id"`
	EventType string    `bson:"event_type"`
	Timestamp time.Time `bson:"timestamp"`
	Viewed    bool      `bson:"viewed"`
	Deleted   bool      `bson:"deleted"`
}
