package api

import (
	"fmt"
	"regexp"

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
	var params types.User
	if err := ctx.BodyParser(&params); err != nil {
		return fiber.ErrBadRequest
	}

	user, errors := validateUser(&params)
	if errors != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	insertedUser, err := h.store.Insert(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

//TODO: get user by id, delete user, update user

// inner
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return emailRegex.MatchString(e)
}
func validateUser(params *types.User) (*types.User, map[string]string) {
	errMap := map[string]string{}
	if len(params.FullName) < minFullName || len(params.FullName) > maxFullName {
		errMap["name"] = fmt.Sprintf("name must be between %d and %d", minNameLen, maxNameLen)
	}
	if !isEmailValid(params.Email) {
		errMap["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	if len(errMap) != 0 {
		return nil, errMap
	}

	return &types.User{
		FullName:          params.FullName,
		Email:             params.Email,
		EncryptedPassword: "sisma", //TODO password bcrypt and JWT or somt
		Admin:             false,
	}, nil
}
