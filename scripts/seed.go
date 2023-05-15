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
	latMin = 32.554
	latMax = 32.80

	langMax = 34.96
	langMin = 35.53
)

var store db.Store

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.GetURI()))
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

	admin := addUser(&store, "a", true)

	fmt.Printf("email\t%s\npassword\t%s\n", "a@a.com", "sisma_a")

	fp := addFoodProvider(&store, "Israel", "item", "haifa")
	rider := AddRider(&store, "Israel", "Israel", user.ID, fp.ID)

	fmt.Println("admin token -->", admin.GenerateToken())
	fmt.Printf("\ntoken:\t%s\n\nfpid:\t%s\nRider:\t%s\n", user.GenerateToken(), fp.ID.Hex(), rider.ID.Hex())

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
			Lat:     latMin + rand.Float64()*(latMax-latMin),
			Lng:     langMin + rand.Float64()*(langMax-langMin),
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
		FoodProviderID: fid,
		Destination: types.Location{
			Address: destAddr,
			Lat:     latMin + rand.Float64()*(latMax-latMin),
			Lng:     langMin + rand.Float64()*(langMax-langMin),
		},
		Available: a,
		Location: types.Location{
			Address: addr,
			Lat:     latMin + rand.Float64()*(latMax-latMin),
			Lng:     langMin + rand.Float64()*(langMax-langMin),
		},
	}

	insertedRider, err := store.Rider.Insert(context.TODO(), &r)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRider
}
