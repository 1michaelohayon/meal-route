package api

import (
	"errors"

	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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

func (h *FoodProviderHandler) Add(ctx *fiber.Ctx) error {
	var params types.NewFoodProvider
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	errors := params.Validate()
	if errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	fp := params.CreateFoodProvider()

	insertedFp, err := h.store.Insert(ctx.Context(), fp)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedFp)
}
