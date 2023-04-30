package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Rider struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId         primitive.ObjectID `bson:"UserId,omitempty" json:"UserId,omitempty"`
	FoodProviderID primitive.ObjectID `bson:"foodProvider,omitempty" json:"foodProvider,omitempty"`
	Destination    Location           `bson:"destination,omitempty" json:"destination,omitempty"`
	Available      bool               `bson:"available" json:"available"`
	Location
}
