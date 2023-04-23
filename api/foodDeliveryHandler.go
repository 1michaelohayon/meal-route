package api

import (
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

func NewFoodProviderHandler(fps db.FoodProviderStore) *FoodProviderHandler {
	return &FoodProviderHandler{
		store: fps,
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

func validate(params *types.FoodProvider) (*types.FoodProvider, map[string]string) {
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
func (h *FoodProviderHandler) Post(ctx *fiber.Ctx) error {
	var params types.FoodProvider
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	fp, errors := validate(&params)
	if len(errors) != 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	insertedFp, err := h.store.Insert(ctx.Context(), fp)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedFp)
}
