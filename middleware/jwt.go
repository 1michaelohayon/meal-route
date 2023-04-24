package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/1michaelohayon/meal-route/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthorizeation(uStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.GetReqHeaders()["Authorization"]
		if !ok {
			fmt.Println("token not present in the header")
			return fiber.ErrUnauthorized
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return ctx.Status(401).JSON(map[string]string{"Authorization error": "token expired"})
		}

		userID := claims["id"].(string)
		user, err := uStore.GetById(ctx.Context(), userID)
		if err != nil {
			fmt.Println(err)
			return fiber.ErrUnauthorized
		}
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fiber.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token", err) // TODO change these to logger
		return nil, fiber.ErrUnauthorized
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}
