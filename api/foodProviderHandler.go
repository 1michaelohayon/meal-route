package api

import (
	"errors"
	"fmt"

	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	minNameLen = 3
	maxNameLen = 50
	minAddrLen = 4
	maxAddrLen = 100
)

type FoodProviderHandler struct {
	store db.FoodProviderStore
}

func NewFoodProviderHandler(f db.FoodProviderStore) *FoodProviderHandler {
	return &FoodProviderHandler{
		store: f,
	}
}

func (h *FoodProviderHandler) GetAll(ctx *fiber.Ctx) error {
	foodProviders, err := h.store.Get(ctx.Context())
	if err != nil {
		if err == mongo.ErrNilDocument {
			return ctx.JSON([]types.FoodProvider{})
		}
		return err
	}

	return ctx.JSON(foodProviders)
}

func (h *FoodProviderHandler) GetOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	foodProvider, err := h.store.GetById(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Status(404).JSON(map[string]string{"foodProvider->GetOne error": "not found"})
		}
		return err
	}
	return ctx.JSON(foodProvider)
}

func validateFP(params *types.FoodProvider) (*types.FoodProvider, map[string]string) {
	// TODO move this func to be attached in types
	errMap := map[string]string{}
	if len(params.Name) < minNameLen || len(params.Name) > maxNameLen {
		errMap["name"] = fmt.Sprintf("name must be between %d and %d", minNameLen, maxNameLen)
	}

	if len(params.Location.Address) < minNameLen || len(params.Location.Address) > maxNameLen {
		errMap["address"] = fmt.Sprintf("address must be between %d and %d", minAddrLen, maxAddrLen)
	}

	if len(errMap) != 0 {
		return nil, errMap
	}

	return &types.FoodProvider{
		Name:     params.Name,
		Location: params.Location,
		Menu:     params.Menu,
	}, nil
}
func (h *FoodProviderHandler) Add(ctx *fiber.Ctx) error {
	var params types.FoodProvider
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	fp, errors := validateFP(&params)
	if errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	insertedFp, err := h.store.Insert(ctx.Context(), fp)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedFp)
}
