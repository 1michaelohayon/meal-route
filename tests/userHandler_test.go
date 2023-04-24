package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/1michaelohayon/meal-route/api"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
)

var (
	_validUser = types.NewUser{
		Email:    "abc@gmail.com",
		FullName: "Valid User",
		Password: "sisma123",
	}
	_invalidUser = types.NewUser{
		Email:    "abc@gmaimadad22",
		FullName: "Vr",
		Password: "si",
	}
)

func TestAddValidUser(t *testing.T) {
	var (
		tdb        = setup(t)
		app        = fiber.New()
		userHanlde = api.NewUserHandler(tdb.Store.User)
		have       types.User
	)
	defer tdb.teardown(t)
	app.Post("/", userHanlde.Add)
	b, _ := json.Marshal(_validUser)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("\nexpected:\t200 status but have %s\n", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&have)
	if len(have.ID) == 0 {
		t.Errorf("\nexpected:\tuser ID to be set\n")
	}
	if len(have.EncryptedPassword) > 0 {
		t.Errorf("\nexpected:\tencryptedpassword not to be included in the json resp\n")
	}
	if have.Email != _validUser.Email {
		t.Errorf("\nexpected\tEmail %s but have %s\n", _validUser.Email, have.Email)
	}
	if have.FullName != _validUser.FullName {
		t.Errorf("\nexpected:\tFullName %s but have %s", _validUser.FullName, have.FullName)
	}
}

func TestAddInvalidUser(t *testing.T) {
	var (
		tdb        = setup(t)
		app        = fiber.New()
		userHanlde = api.NewUserHandler(tdb.Store.User)
		have       types.User
	)
	defer tdb.teardown(t)
	app.Post("/", userHanlde.Add)
	b, _ := json.Marshal(_invalidUser)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("\nexpected:\t400 status but have %s\n", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&have)
	if len(have.Email) != 0 {
		t.Errorf("\nexpected:\tempty but have %s\n", have.Email)
	}
	if len(have.FullName) != 0 {
		t.Errorf("\nexpected:\tempty but have %s\n", have.FullName)
	}
}
