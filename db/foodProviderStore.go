package db

import (
	"context"

	"github.com/1michaelohayon/meal-route/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodProviderStore interface {
	Get(context.Context) ([]*types.FoodProvider, error)
	Insert(context.Context, *types.FoodProvider) (*types.FoodProvider, error)
}
type MongoFoodProvider struct {
	client *mongo.Client
	coll   *mongo.Collection

	FoodProviderStore
}

func NewFoodProviderwMongoStore(cl *mongo.Client) *MongoFoodProvider {
	return &MongoFoodProvider{
		client: cl,
		coll:   cl.Database(DBNAME).Collection("food-providers"),
	}
}

func (s *MongoFoodProvider) Insert(ctx context.Context, fp *types.FoodProvider) (*types.FoodProvider, error) {
	stored, err := s.coll.InsertOne(ctx, fp)
	if err != nil {
		return nil, err
	}
	fp.ID = stored.InsertedID.(primitive.ObjectID).Hex()
	return fp, nil
}

func (s *MongoFoodProvider) Get(ctx context.Context) ([]*types.FoodProvider, error) {
	curs, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var providers []*types.FoodProvider
	if err = curs.All(ctx, &providers); err != nil {
		return nil, err
	}

	return providers, err
}
