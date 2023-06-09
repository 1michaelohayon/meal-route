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
	GetById(ctx context.Context, id string) (*types.Rider, error)
	UpdateLocation(ctx context.Context, pos types.Location, id string) error
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

func (s *MongoRiderStore) GetById(ctx context.Context, id string) (*types.Rider, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var fp types.Rider
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&fp); err != nil {
		return nil, err
	}
	return &fp, nil
}

func (s *MongoRiderStore) UpdateLocation(ctx context.Context, pos types.Location, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := map[string]primitive.ObjectID{"_id": oid}
	update := map[string]bson.M{"$set": {"location": pos}}
	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}
