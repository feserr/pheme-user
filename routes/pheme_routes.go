package routes

import (
	"github.com/feserr/pheme-user/controllers"
	"github.com/gofiber/fiber/v2"
)

// PhemeSetup setups the pheme endpoints.
func PhemeSetup(app *fiber.App) {
	app.Get("/api/v1/pheme", controllers.GetAllPhemes)
	app.Get("/api/v1/pheme/mine", controllers.GetMyPhemes)
	app.Get("/api/v1/pheme/user/:id<int>", controllers.GetUserPhemes)
	app.Get("/api/v1/pheme/:id<int>", controllers.GetPheme)
	app.Post("/api/v1/pheme", controllers.PostPheme)
	app.Delete("/api/v1/pheme/:id<int>", controllers.DeletePheme)
	app.Put("/api/v1/pheme/:id<int>", controllers.UpdatePheme)
}
