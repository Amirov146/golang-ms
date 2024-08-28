package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"golang-ms/src/config"
	"golang-ms/src/config/models"
	"log"
)

type Main struct {
	token     *models.PasetoAuth
	config    *config.Config
	routerApi *fiber.App
}

func newMain(c *config.Config, routerApi *fiber.App) (*Main, error) {

	pasetoToken, err := models.NewPaseto([]byte(c.Token.TokenKey))
	if err != nil {
		return nil, err
	}

	app := &Main{
		token:     pasetoToken,
		routerApi: routerApi,
		config:    c,
	}

	//app.SetApi()

	return app, nil
}

func Start() string {
	return ""
}

func (a *Main) Start() error {
	return a.routerApi.Listen(a.config.Token.Address)
}

func main() {
	c := config.LoadConfig()

	routerApi := fiber.New()
	app, err := newMain(c, routerApi)
	if err != nil {
		panic(fmt.Errorf("create app error:", err))
	}
	//app.SetApi()
	log.Fatal(app.Start())

}
