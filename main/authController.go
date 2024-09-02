package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-ms/src/models"
	"golang-ms/src/service"
	"io/ioutil"
	"net/http"
)

func (a *Main) SetApi() {

	//public routes
	a.routerApi.Post("/login", a.Login)
	a.routerApi.Post("/register", service.CreateUser)

	a.routerApi.Get("/account", func(c *fiber.Ctx) error {

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

	// protected routes
	protectedApi := a.routerApi.Group("/api", a.CheckAuth())
	{
		protectedApi.Get("/all", a.RoleBasedAuth("ADMIN"), func(c *fiber.Ctx) error {
			resp, err := http.Get("http://localhost:8082/all")
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return c.Send(body)
		})
	}
}
