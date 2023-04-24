package api

import (
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderHandler struct {
	store db.RiderStore
}

func NewRiderHandler(s db.RiderStore) *RiderHandler {
	return &RiderHandler{
		store: s,
	}
}

func (h *RiderHandler) GetAll(ctx *fiber.Ctx) error {
	riders, err := h.store.Get(ctx.Context())
	if err != nil {
		if err == mongo.ErrNilDocument {
			return ctx.JSON([]types.Rider{})
		}
		return err
	}
	return ctx.JSON(riders)
}

func (h *RiderHandler) Add(ctx *fiber.Ctx) error {
	var params types.Rider
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}

	//TODO VALIDATION

	insertedRider, err := h.store.Insert(ctx.Context(), &types.Rider{})
	if err != nil {
		return err
	}
	return ctx.JSON(insertedRider)
}
