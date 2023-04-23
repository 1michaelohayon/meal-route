package types

type Location struct {
	Address string  `bson:"address" json:"address"`
	Lat     float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lang    float64 `bson:"lang,omitempty" json:"lang,omitempty"`
}

type FoodProvider struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Menu     []string `bson:"menu" json:"menu"`
	Location Location `bson:"location" json:"location"`
}
