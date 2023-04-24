package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minNameLen = 3
	maxNameLen = 50
	minAddrLen = 4
	maxAddrLen = 100
)

type Location struct {
	Address string  `bson:"address,omitempty" json:"address,omitempty"`
	Lat     float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lang    float64 `bson:"lang,omitempty" json:"lang,omitempty"`
}

type FoodProvider struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Menu     []string           `bson:"menu" json:"menu"` // TODO expend to struct
	Location `bson:"location" json:"location"`
	Riders   []primitive.ObjectID `bson:"riders" json:"riders"`
}
type NewFoodProvider struct {
	Name     string   `bson:"name" json:"name"`
	Menu     []string `bson:"menu" json:"menu"`
	Location Location `bson:"location" json:"location"`
}

func (params *NewFoodProvider) Validate() map[string]string {
	errMap := map[string]string{}
	if len(params.Name) < minNameLen || len(params.Name) > maxNameLen {
		errMap["name error"] = fmt.Sprintf("name must be between %d and %d", minNameLen, maxNameLen)
	}
	if len(params.Location.Address) < minNameLen || len(params.Location.Address) > maxNameLen {
		errMap["address error"] = fmt.Sprintf("address must be between %d and %d", minAddrLen, maxAddrLen)
	}

	if len(errMap) != 0 {
		return errMap
	}
	return nil
}
func (params *NewFoodProvider) CreateFoodProvider() *FoodProvider {
	if params.Menu == nil {
		params.Menu = []string{}
	}

	return &FoodProvider{
		Name:     params.Name,
		Location: params.Location,
		Menu:     params.Menu,
		Riders:   []primitive.ObjectID{},
	}
}
