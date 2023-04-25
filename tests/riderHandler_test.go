package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
)

var (
	_rider = types.Rider{
		Available: true,
		Location: types.Location{
			Address: "haifa",
			Lat:     3.425235,
			Lang:    1.2,
		},
		Destination: types.Location{
			Address: "haifa",
			Lat:     2.425235,
			Lang:    1.5,
		},
	}
)

func TestAddRider(t *testing.T) {
	var (
		tdb                = setup(t)
		app                = fiber.New()
		riderHandle        = api.NewRiderHandler(*tdb.Store)
		foodProviderHandle = api.NewFoodProviderHandler(tdb.Fp)
		have               types.Rider
		fp                 types.FoodProvider
	)
	defer tdb.teardown(t)
	//adding foodProvider to post rider to
	app.Post("/", foodProviderHandle.Add)
	app.Get("/:id", foodProviderHandle.GetOne)
	app.Post("/:foodPorviderID", riderHandle.Add)

	b, _ := json.Marshal(_validFoodProvider)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&fp)

	//adding rider
	b, _ = json.Marshal(_rider)
	req = httptest.NewRequest("POST", "/"+fp.ID.Hex(), bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err = app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("\nexpected: 200 status but have %s\n", resp.Status)
	}
	json.NewDecoder(resp.Body).Decode(&have)
	if err != nil {
		t.Fatal(err)
	}

	if have.FoodProviderID.Hex() != fp.ID.Hex() {
		t.Errorf("\nexpected: rider's food provider id to be %s but have %s\n", fp.ID.Hex(), have.FoodProviderID.Hex())
	}
	//test if the foodProvider was updated
	b, _ = json.Marshal(_validFoodProvider)
	req = httptest.NewRequest("GET", "/"+fp.ID.Hex(), bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&fp)

	if fp.Riders[0].Hex() != have.ID.Hex() {
		t.Errorf("\nexpected: foodProvider.riders list to be updated with the added rider IDs\n")
	}
}

func TestGetAllRider(t *testing.T) {
	var (
		tdb         = setup(t)
		app         = fiber.New()
		riderHandle = api.NewRiderHandler(*tdb.Store)
		have        types.Rider
	)
	defer tdb.teardown(t)
	app.Get("/:foodPorviderID/:id", riderHandle.GetOne)
	app.Post("/:foodPorviderID", riderHandle.Add)

	insertedFp, err := tdb.Store.Fp.Insert(context.TODO(), _validFoodProvider.CreateFoodProvider())
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(_rider)
	req := httptest.NewRequest("POST", "/"+insertedFp.ID.Hex(), bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&have)

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", insertedFp.ID.Hex(), have.ID.Hex()), nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("\nexpected: 200 status but have %s\n", resp.Status)
	}
	var respB types.Rider
	json.NewDecoder(resp.Body).Decode(&respB)

	if respB.ID != have.ID {
		t.Errorf("\nexpected: to find the same rider %s but have %s\n", have.ID.Hex(), respB.ID.Hex())
	}
}
