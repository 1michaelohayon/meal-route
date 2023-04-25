package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	latMin = 29.4965
	latMax = 33.2778

	langMax = 34.2677
	langMin = 35.8961
)

var store db.Store

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store = db.Store{
		User:  db.NewMongoUserStore(client),
		Fp:    db.NewMongoFoodProviderStore(client),
		Rider: db.NewMongoRiderStore(client),
	}
}

func main() {
	for i := 0; i < 3; i++ {
		fp := addFoodProvider(&store, "bamba falafel"+fmt.Sprint(i), "chips", "haifa")

		maxRiders := rand.Intn(4) + 1
		for j := 0; j < maxRiders; j++ {
			u := addUser(&store, "test_user"+fmt.Sprint(i), false)

			r := AddRider(&store, "Israel", "krayot", u.ID, fp.ID)
			if err := store.Fp.PutRider(context.TODO(), fp.ID, r.ID); err != nil {
				log.Fatal(err)
			}
		}
	}

	user := addUser(&store, "Eli_Copter", false)

	admin := addUser(&store, "admin", true)

	fmt.Println("admin token -->", admin.GenerateToken())
	fmt.Printf("\n%s --> %s\n", user.Email, user.GenerateToken())

}

func addUser(store *db.Store, fn string, admin bool) *types.User {
	params := types.NewUser{
		Email:    fmt.Sprintf("%s@%s.com", fn, fn),
		FullName: fn,
		Password: fmt.Sprintf("sisma_%s", fn),
	}

	user, err := params.CreateUser()
	if err != nil {
		log.Fatal(err)
	}

	user.Admin = admin

	insertedUser, err := store.User.Insert(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func addFoodProvider(store *db.Store, name, menu, addr string) *types.FoodProvider {

	fp := types.FoodProvider{
		Name: name,
		Menu: []string{
			fmt.Sprint(menu, 'a'),
			fmt.Sprint(menu, 'b'),
			fmt.Sprint(menu, 'c'),
		},
		Riders: []primitive.ObjectID{},
		Location: types.Location{
			Address: addr,
			Lat:     latMin + rand.Float64()*(latMax-langMin)*-1,
			Lang:    langMin + rand.Float64()*(langMax-langMin)*-1,
		},
	}

	insertedFp, err := store.Fp.Insert(context.TODO(), &fp)
	if err != nil {
		log.Fatal(err)
	}

	return insertedFp
}

func AddRider(store *db.Store, addr, destAddr string, uid, fid primitive.ObjectID) *types.Rider {
	a := true
	if rand.Intn(2) == 0 {
		a = false
	}

	r := types.Rider{
		UserId:         uid,
		FoodProviderID: uid,
		Destination: types.Location{
			Address: destAddr,
			Lat:     latMin + rand.Float64()*(latMax-langMin)*-1,
			Lang:    langMin + rand.Float64()*(langMax-langMin)*-1,
		},
		Available: a,
		Location: types.Location{
			Address: addr,
			Lat:     latMin + rand.Float64()*(latMax-langMin)*-1,
			Lang:    langMin + rand.Float64()*(langMax-langMin)*-1,
		},
	}

	insertedRider, err := store.Rider.Insert(context.TODO(), &r)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRider
}
