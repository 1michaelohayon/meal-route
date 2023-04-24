package db

const (
	DBNAME     = "food-route"
	DBURI      = "mongodb://localhost:27017"
	TestDBNAME = "TEST_food-route"
)

type Store struct {
	User  UserStore
	Fp    FoodProviderStore
	Rider RiderStore
}
