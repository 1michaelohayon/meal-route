package middleware

import (
	"github.com/1michaelohayon/meal-route/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(ctx *fiber.Ctx) error {
	user, ok := ctx.Context().UserValue("user").(*types.User)

	if !ok || !user.Admin {
		return ctx.Status(401).JSON(map[string]string{"admin error": "unauthorized"})
	}
	return ctx.Next()
}
