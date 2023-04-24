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
	listenAddr         = flag.String("listenAddr", ":5000", "Api PORT")
	app                = fiber.New()
	apiRoute           = app.Group("/api")
	foodProviderHandle *api.FoodProviderHandler
	userHandle         *api.UserHandler
	riderHandle        *api.RiderHandler
)

func init() {
	flag.Parse()
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	store := db.Store{
		User:  db.NewMongoUserStore(client),
		Fp:    db.NewMongoFoodProviderStore(client),
		Rider: db.NewMongoRiderStore(client),
	}

	userHandle = api.NewUserHandler(store.User)
	foodProviderHandle = api.NewFoodProviderHandler(store.Fp)
	riderHandle = api.NewRiderHandler(store)
}

func main() {
	/* food-provider routes */
	apiRoute.Get("/foodprovider", foodProviderHandle.GetAll)
	apiRoute.Get("/foodprovider/:id", foodProviderHandle.GetOne)
	apiRoute.Post("/foodprovider", foodProviderHandle.Add)
	/* rider routes */
	apiRoute.Post("/foodprovider/:foodPorviderID/", riderHandle.Add)
	apiRoute.Get("/foodprovider/:foodPorviderID/riders", riderHandle.GetAll)

	/* user routes  */
	apiRoute.Get("/user", userHandle.GetAll)
	apiRoute.Post("/user", userHandle.Add)

	/* admin routes */
	app.Listen(*listenAddr)

}
