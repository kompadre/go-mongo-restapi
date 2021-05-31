package mongodb

import (
	"context"
	"go-mongo-restapi/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoStorageHandle struct {
	Collection *mongo.Collection
	Ctx *context.Context
}

func (h *MongoStorageHandle) RetrieveProducts(maxprice int64) *model.Products {
	var limit int64 = 10
	opts := options.FindOptions{
		Limit: &limit,
	}
	cond := bson.M{}
	if maxprice > 0 {
		cond = bson.M{"original_price": bson.M{"$lt": maxprice}}
	}
	cursor, err := h.Collection.Find(context.Background(), cond, &opts)
	if err != nil {
		log.Panicf("Error retrieving products from mongodb: %v", err)
		return nil
	}
	var products model.Products
	cursor.All(context.Background(), &products.Data)
	return &products
}
