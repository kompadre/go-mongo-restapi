package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-mongo-restapi/internal/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func main() {
	products := product.Retrieve()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongohost := os.Getenv("MONGODB_HOST")
	if (mongohost == "") { // This is for local, non Docker testing and/or debugging.
		mongohost = "localhost"
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", mongohost)))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("testing").Collection("products")
	for _, v := range products.Data {
		res, err := collection.InsertOne(ctx, v)
		if err != nil {
			fmt.Println(err)
		} else {
			id := res.InsertedID
			fmt.Println(id)
		}
	}
	q := &bson.M{"sku":"000002"}
	p := product.Product{}
	e := collection.FindOne(ctx, q).Decode(&p)
	//p.Price.Currency = "EUR"
	fmt.Println(e)
	fmt.Println(p)
	a, err := json.Marshal(p)
	fmt.Println(string(a))
	/*
	res, err := collection.InsertOne(ctx, products.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := res.InsertedID
	*
	fmt.Println(id)
	 */


	app := fiber.New()
	app.Use(cors.New())

	app.Get("/products", func (c *fiber.Ctx) error {
		c.JSON(p)
		return nil
	})

	app.Listen(":8080")
}
