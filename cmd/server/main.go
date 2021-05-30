package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-mongo-restapi/internal/model"
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
	mongo  *mongo.Client
	logger *log.Logger
}

func (server *Server) routes() {
	server.app.Get("/products", func(c *fiber.Ctx) error {
		ps := model.RetrieveFromJson()
		return c.JSON(ps)
	})
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
	mongoclient, _ := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))

	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	server.mongo = mongoclient

	app := fiber.New()
	app.Use(cors.New())
	server.app = app
	server.routes()
	server.app.Listen(":" + strconv.Itoa(Config.http.port))

	return nil
}

func main() {
	server := Server{}
	if err := server.run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
