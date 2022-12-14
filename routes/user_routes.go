package routes

import (
	authModels "github.com/feserr/pheme-auth/controllers"
	"github.com/feserr/pheme-user/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserSetup(app *fiber.App) {
	app.Get("/api/v1/user", authModels.User)
	app.Get("/api/v1/user/:name<string>", controllers.GetUsersByName)
	app.Put("/api/v1/user/friend/:id<int>", controllers.AddFriend)
	app.Put("/api/v1/user/follower/:id<int>", controllers.AddFollower)
	app.Delete("/api/v1/user/friend/:id<int>", controllers.DeleteFriend)
	app.Delete("/api/v1/user/follower/:id<int>", controllers.DeleteFollower)
}
