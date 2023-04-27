package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
)

var (
	_validFoodProvider = types.NewFoodProvider{
		Name: "MakTchina",
		Menu: []string{
			"corn with water",
			"potato with canola oil",
			"fish for cats",
		},
		Location: types.Location{
			Address: "negev",
			Lat:     30.015463740207814,
			Lng:     36.20807917170951,
		},
	}
	_InvalidFoodProvider = types.NewFoodProvider{
		Name: "1",
		Location: types.Location{
			Address: "2",
		},
	}
)

func TestAddValidFoodProvider(t *testing.T) {
	var (
		tdb                = setup(t)
		app                = fiber.New()
		foodProviderHandle = api.NewFoodProviderHandler(tdb.Store.Fp)
		have               types.FoodProvider
	)
	defer tdb.teardown(t)
	app.Post("/", foodProviderHandle.Add)
	b, _ := json.Marshal(_validFoodProvider)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("\nexpected: 200 status but have %s\n", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&have)

	if len(have.ID) == 0 {
		t.Errorf("\nexpected: user ID to be set\n")
	}
	if have.Name != _validFoodProvider.Name {
		t.Errorf("\nexpected: Name %s but have %s", _validFoodProvider.Name, have.Name)
	}

	if _, err = tdb.Store.Fp.GetById(context.TODO(), have.ID.Hex()); err != nil {
		log.Fatal(err)
	}
}

func TestAddInvalidFoodProvider(t *testing.T) {
	var (
		tdb                = setup(t)
		app                = fiber.New()
		foodProviderHandle = api.NewFoodProviderHandler(tdb.Store.Fp)
		have               types.FoodProvider
	)
	defer tdb.teardown(t)
	app.Post("/", foodProviderHandle.Add)
	b, _ := json.Marshal(_InvalidFoodProvider)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("\nexpected: 400 status but have %s\n", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&have)

	foodProviders, err := tdb.Fp.Get(context.TODO())
	if err != nil {
		t.Error(err)
	}
	if len(have.Name) > 0 {
		t.Errorf("\nexpected: Name empty but have %s", have.Name)
	}
	if len(foodProviders) != 0 {
		t.Errorf("\nexpected: 0 foodProviders but have %d", len(foodProviders))
	}
}

func TestGetAllFoodProvider(t *testing.T) {
	var (
		tdb                = setup(t)
		app                = fiber.New()
		foodProviderHandle = api.NewFoodProviderHandler(tdb.Store.Fp)
		have               []types.FoodProvider
	)
	defer tdb.teardown(t)

	app.Get("/", foodProviderHandle.GetAll)

	first, err := tdb.Store.Fp.Insert(context.TODO(), _validFoodProvider.CreateFoodProvider())
	if err != nil {
		t.Error(err)
	}
	second, err := tdb.Store.Fp.Insert(context.TODO(), _validFoodProvider.CreateFoodProvider())
	if err != nil {
		t.Error(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(resp.Body).Decode(&have)

	if len(have) != 2 {
		t.Errorf("\nexpected: to have 2 foodProviders but have %d\n", len(have))
	}

	list := []types.FoodProvider{
		*first, *second,
	}

	if !reflect.DeepEqual(list, have) {
		t.Errorf("\nexpected: resp to be %+v but have %+v\n", list, have)
	}
}

func TestGetOneFoodProvider(t *testing.T) {
	var (
		tdb                = setup(t)
		app                = fiber.New()
		foodProviderHandle = api.NewFoodProviderHandler(tdb.Store.Fp)
		have               types.FoodProvider
	)
	defer tdb.teardown(t)

	app.Get("/:id", foodProviderHandle.GetOne)

	first, err := tdb.Store.Fp.Insert(context.TODO(), _validFoodProvider.CreateFoodProvider())
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("GET", "/"+first.ID.Hex(), nil)
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	json.NewDecoder(resp.Body).Decode(&have)

	if !reflect.DeepEqual(*first, have) {
		t.Errorf("\nexpected: resp to be %+v but have %+v\n", *first, have)
	}
}
