package db

import (
	"context"

	"github.com/1michaelohayon/meal-route/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	Get(context.Context) ([]*types.User, error)
	Insert(context.Context, *types.User) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	UserStore
}

func NewMongoUserStore(cl *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: cl,
		coll:   cl.Database(DBNAME).Collection("users"),
	}
}

func (s *MongoUserStore) Insert(ctx context.Context, user *types.User) (*types.User, error) {
	stored, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = stored.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (s *MongoUserStore) Get(ctx context.Context) ([]*types.User, error) {
	curs, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err = curs.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
