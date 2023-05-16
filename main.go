package main

import (
	"context"
	"flag"
	"log"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	//"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	listenAddr = flag.String("listenAddr", ":5000", "server PORT")
	app        = fiber.New()

	apiRoute  fiber.Router
	userRoute fiber.Router
	fpRoute   fiber.Router

	foodProviderHandle *api.FoodProviderHandler
	userHandle         *api.UserHandler
	riderHandle        *api.RiderHandler
	authHandle         *api.AuthHandler
)

func init() {
	uri := db.GetURI()
	flag.Parse()
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri))
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

	app.Use(logger.New())
	app.Use(cors.New())

	apiRoute = app.Group("/api")
	fpRoute = apiRoute.Group("/foodprovider", middleware.JWTAuthorizeation(store.User))
	userRoute = apiRoute.Group("/user", middleware.JWTAuthorizeation(store.User))
}

func main() {

	/* auth */
	apiRoute.Post("/auth", authHandle.Authenticate)
	apiRoute.Post("/signup", userHandle.Add)

	/* food-provider routes */
	fpRoute.Get("/", foodProviderHandle.GetAll)
	fpRoute.Get("/:id", foodProviderHandle.GetOne)
	fpRoute.Post("/", foodProviderHandle.Add)

	/* rider routes */
	fpRoute.Post("/:foodPorviderID/", riderHandle.Add)
	fpRoute.Get("/:foodPorviderID/riders", riderHandle.GetAll)
	fpRoute.Get("/:foodPorviderID/:id", riderHandle.GetOne)
	fpRoute.Put("/:foodPorviderID/:id", riderHandle.SetPosition)

	/* user routes  */
	userRoute.Get("/", middleware.AdminAuth, userHandle.GetAll) /* admin route */
	userRoute.Get("/:id", userHandle.GetOne)

	app.Listen(*listenAddr)
}

/* TODO:
middleware for admin
tests
finetune errors
*/
