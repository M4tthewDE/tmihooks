package db

import (
	"context"
	"time"

	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/m4tthewde/tmihooks/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHandler struct {
	Config *config.Config
}

func (dbHandler *DatabaseHandler) Insert(webhook *structs.Webhook) primitive.ObjectID {
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

	return res.InsertedID.(primitive.ObjectID)
}

func (dbHandler *DatabaseHandler) Delete(webhook *structs.Webhook) int64 {
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

	res, err := collection.DeleteOne(ctx, doc)
	if err != nil {
		panic(err)
	}

	return res.DeletedCount
}

func (dbHandler *DatabaseHandler) DeleteByID(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: dbHandler.Config.Database.Password,
		Password: dbHandler.Config.Database.Password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		return 0, err
	}

	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	collection := client.Database("dev").Collection("webhooks")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": hexID})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (dbHandler *DatabaseHandler) SetConfirmed(id string) int64 {
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

	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	resp, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": hexID},
		bson.D{
			{Key: "$set", Value: bson.M{"status": "confirmed"}},
		},
	)
	if err != nil {
		panic(err)
	}

	return resp.ModifiedCount
}

func (dbHandler *DatabaseHandler) GetWebhook(id string) (*structs.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: dbHandler.Config.Database.Password,
		Password: dbHandler.Config.Database.Password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		return nil, err
	}

	collection := client.Database("dev").Collection("webhooks")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cur, err := collection.Find(
		ctx,
		bson.M{"_id": hexID},
	)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var webhook structs.Webhook
	for cur.Next(context.Background()) {
		err = cur.Decode(&webhook)
		if err != nil {
			return nil, err
		}
	}

	return &webhook, nil
}

func (dbHandler *DatabaseHandler) GetAllWebhooks() ([]*structs.Webhook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username: dbHandler.Config.Database.Password,
		Password: dbHandler.Config.Database.Password,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	if err != nil {
		return nil, err
	}

	collection := client.Database("dev").Collection("webhooks")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(
		ctx,
		bson.D{},
	)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var webhooks []*structs.Webhook
	var webhook structs.Webhook
	for cur.Next(context.Background()) {
		err = cur.Decode(&webhook)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, &webhook)
	}

	return webhooks, nil
}
