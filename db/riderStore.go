package db

import (
	"context"

	"github.com/1michaelohayon/meal-route/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderStore interface {
	Get(context.Context) ([]*types.Rider, error)
	Insert(context.Context, *types.Rider) (*types.Rider, error)
}
type MongoRiderStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	RiderStore
}

func NewMongoRiderStore(cl *mongo.Client) *MongoRiderStore {
	return &MongoRiderStore{
		client: cl,
		coll:   cl.Database(DBNAME).Collection("riders"),
	}
}

func (s *MongoRiderStore) Insert(ctx context.Context, rider *types.Rider) (*types.Rider, error) {
	stored, err := s.coll.InsertOne(ctx, rider)
	if err != nil {
		return nil, err
	}
	rider.ID = stored.InsertedID.(primitive.ObjectID)
	return rider, nil

}

func (s *MongoRiderStore) Get(ctx context.Context) ([]*types.Rider, error) {
	curs, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var riders []*types.Rider
	if err = curs.All(ctx, &riders); err != nil {
		return nil, err
	}

	return riders, nil
}
