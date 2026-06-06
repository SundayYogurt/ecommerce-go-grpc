package middlewares

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sundayyogurt/user_service/internal/dto"
)

func AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// grab authorization header
		authHeader := ctx.Get("Authorization")

		// verify token
		user, err := VerifyToken(authHeader)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// if token is valid, proceed to the next handler
		// assign the decoded suer to the context
		ctx.Locals("userID", user.UserID)
		ctx.Locals("user", user)
		return ctx.Next()
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

func VerifyToken(tokenString string) (dto.AuthResponse, error) {
	tokenArr := strings.Split(tokenString, " ")
	if len(tokenArr) != 2 || tokenArr[0] != "Bearer" {
		return dto.AuthResponse{}, errors.New("invalid token format")
	}

	tokenStr := tokenArr[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return dto.AuthResponse{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return dto.AuthResponse{}, errors.New("invalid token claims")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return dto.AuthResponse{}, errors.New("token expired")
	}

	return dto.AuthResponse{
		UserID: int(claims["user_id"].(float64)),
		Email:  claims["email"].(string),
		Exp:    claims["exp"].(float64),
	}, nil
}
