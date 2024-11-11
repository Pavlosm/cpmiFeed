package main

import (
	"context"
	"cpmiFeed/pkg/db"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	mDb := db.NewClient().Database(db.GetDatabase())
	defer mDb.Client().Disconnect(context.TODO())
	episodesCollection := mDb.Collection("Event")
	var wg sync.WaitGroup

	eventStream, err := episodesCollection.Watch(context.TODO(), mongo.Pipeline{})

	if err != nil {
		panic(err)
	}

	routineCtx, _ := context.WithCancel(context.Background())
	go iterateChangeStream(routineCtx, &wg, eventStream)
	wg.Wait()
}

func iterateChangeStream(routineCtx context.Context, wg *sync.WaitGroup, stream *mongo.ChangeStream) {
	wg.Add(1)
	defer stream.Close(routineCtx)
	defer wg.Done()
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", data)
	}
}
