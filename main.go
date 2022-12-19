package main

import (
	"fmt"
	"os"

	_ "github.com/feserr/pheme-user/docs"
	"github.com/feserr/pheme-user/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Pheme users
// @version 1.0
// @description Pheme users service.

// @contact.name Elias Serrano
// @contact.email feserr3@gmail.com

// @license.name BSD 2-Clause License
// @license.url http://opensource.org/licenses/BSD-2-Clause

// @BasePath /api/
func main() {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Pheme users v1.0.0",
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	err := app.Listen(fmt.Sprintf("%v:%v", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
	if err != nil {
		panic(err.Error())
	}
}
