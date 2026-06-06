package middlewares

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// grab authorization header
		// verify token token
		// if to ken is valid, proceed to the next handler
		// assign the decoded suer to the context
		return c.Next()
	}
}

func GenerateToken(userID int, email string) (string, error) {
	if userID == 0 || email == "" {
		return "", errors.New("userID and email are required")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
