package main

import (
	"context"
	"flag"
	"log"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	listenAddr = flag.String("listenAddr", ":5000", "server PORT")
	app        = fiber.New()

	apiRoute fiber.Router

	foodProviderHandle *api.FoodProviderHandler
	userHandle         *api.UserHandler
	riderHandle        *api.RiderHandler
	authHandle         *api.AuthHandler
)

func init() {
	flag.Parse()
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	store := db.Store{
		User:  db.NewMongoUserStore(client),
		Fp:    db.NewMongoFoodProviderStore(client),
		Rider: db.NewMongoRiderStore(client),
	}

	authHandle = api.NewAuthHandler(store.User)
	userHandle = api.NewUserHandler(store.User)
	foodProviderHandle = api.NewFoodProviderHandler(store.Fp)
	riderHandle = api.NewRiderHandler(store)

	apiRoute = app.Group("/api", middleware.JWTAuthorizeation(store.User))
}

func main() {
	/* auth */
	app.Post("/auth", authHandle.Authenticate)

	/* food-provider routes */
	apiRoute.Get("/foodprovider", foodProviderHandle.GetAll)
	apiRoute.Get("/foodprovider/:id", foodProviderHandle.GetOne)
	apiRoute.Post("/foodprovider", foodProviderHandle.Add)

	/* rider routes */
	apiRoute.Post("/foodprovider/:foodPorviderID/", riderHandle.Add)
	apiRoute.Get("/foodprovider/:foodPorviderID/riders", riderHandle.GetAll)
	apiRoute.Get("/foodprovider/:foodPorviderID/:id", riderHandle.GetOne)

	/* user routes  */
	apiRoute.Get("/user", userHandle.GetAll)
	apiRoute.Get("/user/:id", userHandle.GetOne)

	apiRoute.Post("/user", userHandle.Add) // TODO change group, can't create new user because it's grouped with the middleware

	/* admin routes */

	app.Listen(*listenAddr)
}

/* TODO:
middleware for admin
tests
finetune errors
*/
