package main

import (
	"context"
	"flag"
	"log"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	listenAddr          = flag.String("listenAddr", ":5000", "The listen address of the API server")
	app                 = fiber.New()
	apiRoute            = app.Group("/api")
	foodProviderHandler *api.FoodProviderHandler
	userHandler         *api.UserHandler
)

func init() {
	flag.Parse()
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	foodProviderStore := db.NewMongoFoodProviderStore(client)
	userStore := db.NewMongoUserStore(client)

	foodProviderHandler = api.NewFoodProviderHandler(foodProviderStore)
	userHandler = api.NewUserHandler(userStore)
}

func main() {
	/* food-provider routes */
	apiRoute.Get("/foodprovider", foodProviderHandler.GetAll)
	apiRoute.Post("/foodprovider", foodProviderHandler.Post)

	/* user routes  */
	apiRoute.Get("/user", userHandler.GetAll)
	apiRoute.Post("/user", userHandler.Post)

	/* admin routes */
	app.Listen(*listenAddr)
}