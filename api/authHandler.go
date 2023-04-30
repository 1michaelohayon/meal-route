package api

import (
	"errors"

	"github.com/1michaelohayon/meal-route/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	uStore db.UserStore
}

func NewAuthHandler(s db.UserStore) *AuthHandler {
	return &AuthHandler{
		uStore: s,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ap *AuthParams) validatePassword(encpw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(ap.Password)) == nil
}

func (h *AuthHandler) Authenticate(ctx *fiber.Ctx) error {
	var params AuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.uStore.GetByEmail(ctx.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Status(404).JSON(map[string]string{"Authenticate error": "user not found"})
		}
		return err
	}

	if !params.validatePassword(user.EncryptedPassword) {
		return ctx.Status(401).JSON(map[string]string{"Authenticate error": "invalid credentials"})
	}

	return ctx.JSON(map[string]string{"token": user.GenerateToken(), "fullName": user.FullName})
}
