package routes

import (
	"pos-services/app/controllers"

	fiber "github.com/gofiber/fiber/v2"
)

func PrivateRoute(a *fiber.App) {
	routeV1 := a.Group("/api/v1")

	routeV1.Post("/user/change-address/:id", controllers.AddUserAddress)
}
