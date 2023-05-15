package db

import (
	"fmt"
	"os"
)

const (
	DBNAME      = "food-route"
	LOCAL_DBURI = "mongodb://localhost:27017"
	PROD_DBURI  = "mongodb://mongodb:27017"
	TestDBNAME  = "TEST_food-route"
)

func GetURI() string {
	env_state := os.Getenv("PROD")
	if env_state == "true" {
		fmt.Println("DB state: Production")
		return PROD_DBURI
	} else {
		fmt.Println("DB state: Local")
		return LOCAL_DBURI
	}
}

type Store struct {
	User  UserStore
	Fp    FoodProviderStore
	Rider RiderStore
}
