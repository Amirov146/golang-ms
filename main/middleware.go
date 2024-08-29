package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-ms/src/models"
	"strings"
)

const (
	authHeader = "Authorization"
	typeBearer = "bearer"
)

var (
	ErrMissingAuthHeader = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

func (a *Main) CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authVal := c.Get(authHeader)

		if len(authVal) == 0 {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": ErrMissingAuthHeader.Error()})
		}

		splitHeader := strings.Fields(authVal)

		if len(splitHeader) < 2 {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": ErrInvalidAuthHeader.Error()})
		}

		authType := strings.ToLower(splitHeader[0])
		if authType != typeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": err.Error()})
		}

		claims, err := a.Token.VerifyToken(splitHeader[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}

func (a *Main) Login(c *fiber.Ctx) error {
	var creds models.Credentials

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	expectedPass, ok := models.Users[creds.Username]
	if !ok || expectedPass != creds.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad password or username"})
	}

	pasetoToken, err := a.Token.NewToken(models.TokenData{
		Subject:  "for user",
		Duration: a.Config.Token.TokenDuration,
		AdditionalClaims: models.AdditionalClaims{
			Name: creds.Username,
			Role: creds.Username,
		},
		Footer: models.Footer{MetaData: "footer for " + creds.Username},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": pasetoToken})
}
