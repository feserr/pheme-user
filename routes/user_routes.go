package routes

import (
	"github.com/feserr/pheme-user/controllers"
	"github.com/gofiber/fiber/v2"
)

// UserSetup setup the user endpoints.
func UserSetup(app *fiber.App) {
	app.Get("/api/v1/users", controllers.GetUsers)
	app.Get("/api/v1/user/:id<int>", controllers.GetUserByID)
	app.Get("/api/v1/user/:name<string>", controllers.GetUsersByName)
	app.Get("/api/v1/user/friend", controllers.GetFriends)
	app.Get("/api/v1/user/follower/:id<int>", controllers.GetFollowers)
	app.Put("/api/v1/user/friend/:id<int>", controllers.AddFriend)
	app.Put("/api/v1/user/follower/:id<int>", controllers.AddFollower)
	app.Delete("/api/v1/user/friend/:id<int>", controllers.DeleteFriend)
	app.Delete("/api/v1/user/follower/:id<int>", controllers.DeleteFollower)
}
