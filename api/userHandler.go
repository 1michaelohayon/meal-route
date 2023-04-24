package api

import (
	"github.com/1michaelohayon/meal-route/db"
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	minFullName = 3
	maxFullName = 75
	minPass     = 7
	maxPass     = 40
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
	insertedUser, err := h.store.Insert(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

//TODO: get user by id, delete user, update user

// inner
