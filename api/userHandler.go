package api

import (
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	uStore db.UserStore
}

func NewUserHandler(s db.UserStore) *UserHandler {
	return &UserHandler{
		uStore: s,
	}
}
func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.uStore.Get(ctx.Context())
	if err != nil {
		if err == mongo.ErrNilDocument {
			return ctx.JSON([]types.User{})
		}
		return err
	}

	return ctx.JSON(users)
}

func (h *UserHandler) Add(ctx *fiber.Ctx) error {
	var params types.NewUser
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}
	if errors := params.Validate(); errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}
	user, err := params.CreateUser()
	if err != nil {
		return err
	}
	insertedUser, err := h.uStore.Insert(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

func (h *UserHandler) GetOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if _, err := h.uStore.GetById(ctx.Context(), id); err != nil {
		return ctx.Status(404).JSON(map[string]string{"User->GetOne error": "not found"})
	}

	user, err := h.uStore.GetById(ctx.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return ctx.JSON(user)
}
