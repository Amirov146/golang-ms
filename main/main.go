package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-ms/src/config"
	"golang-ms/src/models"
	"log"
)

type Main struct {
	Token     *models.PasetoAuth
	Config    *config.Config
	routerApi *fiber.App
}

func newMain(c *config.Config, routerApi *fiber.App) (*Main, error) {

	pasetoToken, err := models.NewPaseto([]byte(c.Token.TokenKey))
	if err != nil {
		return nil, err
	}

	app := &Main{
		Token:     pasetoToken,
		routerApi: routerApi,
		Config:    c,
	}

	app.SetApi()

	return app, nil
}

func (a *Main) Start() error {
	return a.routerApi.Listen(a.Config.Token.Address)
}

func main() {
	c := config.LoadConfig()

	routerApi := fiber.New()
	app, err := newMain(c, routerApi)
	if err != nil {
		panic(fmt.Errorf("create app error:", err))
	}
	app.SetApi()

	//postgresql
	config.ConnectPostgresDB()

	//mongodb
	config.ConnectMongoDB()
	//defer config.DisconnectMongoDB()
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	//<-stop

	log.Fatal(app.Start())
}
