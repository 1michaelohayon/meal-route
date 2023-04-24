package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Address string  `bson:"address,omitempty" json:"address,omitempty"`
	Lat     float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lang    float64 `bson:"lang,omitempty" json:"lang,omitempty"`
}

type FoodProvider struct {
	ID       string               `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Menu     []string             `bson:"menu" json:"menu"`
	Location Location             `bson:"location" json:"location"`
	Riders   []primitive.ObjectID `bson:"riders" json:"riders"`
}

type Rider struct {
	ID             string               `bson:"_id,omitempty" json:"id,omitempty"`
	UserId         primitive.ObjectID   `bson:"UserId,omitempty" json:"UserId,omitempty"`
	FoodProviderID []primitive.ObjectID `bson:"foodProvider,omitempty" json:"foodProvider,omitempty"`
	Destination    Location             `bson:"Destination,omitempty" json:"Destination,omitempty"`

	Location
}
