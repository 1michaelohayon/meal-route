package api

import (
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store db.UserStore
}

func NewUserHandler(s db.UserStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}
func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.store.Get(ctx.Context())
	if err != nil {
		if err == mongo.ErrNilDocument {
			return ctx.JSON([]types.FoodProvider{})
		}
		return err
	}

	return ctx.JSON(users)
}
func (h *UserHandler) Post(ctx *fiber.Ctx) error {
	var params types.User
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}

	//TODO validate

	insertedFp, err := h.store.Insert(ctx.Context(), &params)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedFp)
}
