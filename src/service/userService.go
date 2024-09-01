package service

import (
	"github.com/gofiber/fiber/v2"
	"golang-ms/src/config"
	"golang-ms/src/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GetHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func FindByUsername(username string) (models.User, error) {
	var user models.User

	result := config.DB.Where("username = ?", username).Preload("Roles").First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Please provide the correct input!")
	}

	hashPassword, err := GetHashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something bad happened on the server :(")
	}

	query := `INSERT INTO users (username, first_name, last_name, password, email, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	res := config.DB.Exec(query, user.Username, user.FirstName, user.LastName, hashPassword, user.Email, time.Now())
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something bad happened on the server :(")
	}

	userFound, _ := FindByUsername(user.Username)
	insertRoleQuery := `INSERT INTO users_roles (user_id,role_id) VALUES (?, 1)`
	rs := config.DB.Exec(insertRoleQuery, userFound.ID)
	if rs.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something bad happened on the server :(")
	}

	return c.Status(fiber.StatusOK).SendString("User created successfully")
}
