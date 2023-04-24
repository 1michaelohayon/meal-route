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
	GetById(context.Context, string) (*types.FoodProvider, error)
	PutRider(ctx context.Context, fpId, riderID primitive.ObjectID) error
}
type MongoFoodProviderStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	FoodProviderStore
}

func NewMongoFoodProviderStore(cl *mongo.Client) *MongoFoodProviderStore {
	return &MongoFoodProviderStore{
		client: cl,
		coll:   cl.Database(DBNAME).Collection("food-providers"),
	}
}

func (s *MongoFoodProviderStore) Insert(ctx context.Context, fp *types.FoodProvider) (*types.FoodProvider, error) {
	stored, err := s.coll.InsertOne(ctx, fp)
	if err != nil {
		return nil, err
	}
	fp.ID = stored.InsertedID.(primitive.ObjectID)
	return fp, nil
}

func (s *MongoFoodProviderStore) Get(ctx context.Context) ([]*types.FoodProvider, error) {
	curs, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var providers []*types.FoodProvider
	if err = curs.All(ctx, &providers); err != nil {
		return nil, err
	}

	return providers, nil
}

func (s *MongoFoodProviderStore) GetById(ctx context.Context, id string) (*types.FoodProvider, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var fp types.FoodProvider
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&fp); err != nil {
		return nil, err
	}
	return &fp, nil
}

func (s *MongoFoodProviderStore) PutRider(ctx context.Context, id, riderID primitive.ObjectID) error {
	filter := map[string]primitive.ObjectID{"_id": id}
	update := map[string]bson.M{"$push": {"riders": riderID}}
	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}
