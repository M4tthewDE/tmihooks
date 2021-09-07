package db

import (
	"context"
	"log"
	"time"

	"github.com/m4tthewde/tmihooks/internal/webhook"
	"github.com/m4tthewde/tmihooks/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Insert(webhook webhook.Webhook) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: "test",
		Password: "test",
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		panic(err)
	}

	collection := client.Database("dev").Collection("webhooks")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := util.ToDoc(webhook)
	if err != nil {
		log.Panic(err)
	}
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}
	log.Println(res)
}
