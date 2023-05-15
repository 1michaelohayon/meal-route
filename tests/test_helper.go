package tests

import (
	"context"
	"log"
	"testing"

	"github.com/1michaelohayon/meal-route/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.GetURI()))
	if err != nil {
		log.Fatal(err)
	}

	return &testDb{
		client: client,
		Store: &db.Store{
			User:  db.NewMongoUserStore(client),
			Fp:    db.NewMongoFoodProviderStore(client),
			Rider: db.NewMongoRiderStore(client),
		},
	}
}
