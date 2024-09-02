package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-ms/src/models"
	"golang-ms/src/service"
)

func (a *Main) SetApi() {

	//public routes
	a.routerApi.Post("/login", a.Login)
	a.routerApi.Post("/register", service.CreateUser)

	// protected routes
	protectedApi := a.routerApi.Group("/api", a.CheckAuth())
	{
		//test permissions
		protectedApi.Get("/account", a.RoleBasedAuth("USER"), func(c *fiber.Ctx) error {

			val := c.Locals("claims")

			v, ok := val.(*models.ServiceClaims)

			if !ok {
				return c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": "type conversion error"})
			}

			fmt.Printf("%#v", v)

			owner := fmt.Sprintf("<h3>Account owner - %s</h3>", v.Name)
			role := fmt.Sprintf("<h3>Account role - %s</h3>", v.Role)
			footer := fmt.Sprintf("<h3>Account footer - %s</h3>", v.MetaData)

			return c.SendString(owner + role + footer)
		})

		//users
		protectedApi.Get("/main", a.RoleBasedAuth("USER"), func(c *fiber.Ctx) error {
			return c.Redirect("http://localhost:8082/main")
		})

		//admins
		protectedApi.Get("/admin-panel", a.RoleBasedAuth("ADMIN"), func(c *fiber.Ctx) error {
			return c.Redirect("http://localhost:8082/admin-panel")
		})
	}
}
