package db

import (
	"context"
	"log"
	"time"

	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/util"
	"github.com/m4tthewde/tmihooks/internal/webhook"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHandler struct {
	Config *config.Config
}

func (dbHandler *DatabaseHandler) Insert(webhook webhook.Webhook) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: dbHandler.Config.Database.Password,
		Password: dbHandler.Config.Database.Password,
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
		panic(err)
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}

	log.Println(res)
}
