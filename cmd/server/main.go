package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-mongo-restapi/internal/model"
	"go-mongo-restapi/internal/storage/jsonfile"
	"go-mongo-restapi/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"time"
)

var Config struct {
	mongo struct {
		host string
		port int
	}
	http struct {
		port int
	}
}

type Server struct {
	app    *fiber.App
	mongo  struct {
		client *mongo.Client
		ctx *context.Context
	}
	logger *log.Logger
}

func (server *Server) routes() {
	server.app.Get("/products", func(c *fiber.Ctx) error {
		var ps *model.Products
		if (os.Getenv("TESTING") == "") {
			collection := server.mongo.client.Database("testing").Collection("product")
			log.Printf("Collection: %v", collection)
			ctx, cancel := server.Ctx()
			defer cancel()
			mhandle := mongodb.MongoStorageHandle{Collection: collection, Ctx: ctx}
			q := c.Query("priceLessThan", "0")
			qint, err := strconv.Atoi(q)
			if err != nil {
				qint = 10
			}
			ps = mhandle.RetrieveProducts(int64(qint))
		} else {
			ps = jsonfile.RetrieveProducts()
		}
		ps.ApplyDiscounts(model.CurrentDiscounts())
		fmt.Println(ps)
		return c.JSON(ps)
	})
}

func (server *Server) Ctx() (*context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &ctx, cancel
}

func init() {
	mongohost := os.Getenv("MONGODB_HOST")
	if mongohost == "" { // This is for local, non Docker, testing and/or debugging
		mongohost = "localhost"
	}
	Config.mongo.host = mongohost
	Config.mongo.port = 27017

	httpport, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if httpport == 0 { // This is for local, non Docker, testing and/or debugging
		httpport = 8080
	}
	Config.http.port = httpport
}

func (server *Server) run() error {
	logger := log.New(os.Stderr, "", 1)
	server.logger = logger

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongouri := fmt.Sprintf("mongodb://%s:%d", Config.mongo.host, Config.mongo.port)
	var (
		mongoclient *mongo.Client
		err error
	)

	if os.Getenv("TESTING") == "" {
		// NO time to mock and test this interaction with MongoDB
		// Besides it would still be outside the scope of this application
		mongoclient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
		if err != nil {
			log.Panicf("Error connecting to mongo: %v", err)
		}
		defer func() {
			mongoclient.Disconnect(ctx)
		}()
		server.mongo.client = mongoclient
		server.mongo.ctx = &ctx

		// This, ofcourse, is not part of the products service bur rather we're feeding mongo initial data
		// Normally mongodb's data would be managed outside

		products := jsonfile.RetrieveProducts()
		productsCollection := server.mongo.client.Database("testing").Collection("product")
		productsCollection.DeleteMany(context.Background(), bson.M{})
		_, err = productsCollection.Indexes().CreateOne(context.Background(),mongo.IndexModel{
			Keys:    bson.D{{
				Key:   "sku",
				Value: 1,
			}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			log.Fatalf("Error creating index: %v", err)
		}

		for _, product := range products.Data {
			product.Price.Original = product.OriginalPrice
			product.Price.Final = product.OriginalPrice
			product.Price.Currency = "EUR"
			res, err := productsCollection.InsertOne(*server.mongo.ctx, product)
			if err != nil {
				log.Printf("Failed to insert a product: %v", err)
			} else {
				log.Println(res)
			}
		}
	}

	app := fiber.New()
	app.Use(cors.New()) // I use CORS so that one can access the API from swagger editor for example
	server.app = app
	server.routes()
	err = server.app.Listen(":" + strconv.Itoa(Config.http.port))
	if err != nil {
		log.Panicf("Error staring http server: %v", err)
	}
	return nil
}

func main() {
	server := Server{}
	if err := server.run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
