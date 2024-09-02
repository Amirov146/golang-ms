package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-ms/src/config"
	"golang-ms/src/models"
	"golang-ms/src/service"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	authHeader = "Authorization"
	typeBearer = "bearer"
)

var (
	ErrMissingAuthHeader = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ctx                  = context.Background()
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

func SaveTokenToMongo(token string, username string, expirationTime int64) error {
	ctx := context.Background()

	_, err := config.TokenCollection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"token": token, "expiresAt": expirationTime}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (a *Main) RoleBasedAuth(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("claims")
		claims, ok := val.(*models.ServiceClaims)

		if !ok {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "type conversion error"})
		}

		for _, role := range claims.Role {
			if role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "insufficient permissions"})
	}
}

func (a *Main) Login(c *fiber.Ctx) error {
	var creds models.Credentials

	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userFound, err := service.FindByUsername(creds.Username)
	passError := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(creds.Password))
	if err != nil || passError != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad password or username"})
	}

	var roleNames []string
	for _, role := range userFound.Roles {
		roleNames = append(roleNames, role.Name)
	}

	pasetoToken, err := a.Token.NewToken(models.TokenData{
		Subject:  "for user",
		Duration: a.Config.AccessToken.TokenDuration,
		AdditionalClaims: models.AdditionalClaims{
			Name: userFound.Username,
			Role: roleNames,
		},
		Footer: models.Footer{MetaData: "footer for " + creds.Username},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	expirationTime := time.Now().Add(a.Config.AccessToken.TokenDuration).Unix()
	err = SaveTokenToMongo(pasetoToken, userFound.Username, expirationTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cache saving error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": pasetoToken})
}
