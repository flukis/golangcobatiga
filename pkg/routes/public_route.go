package routes

import (
	"pos-services/app/controllers"

	fiber "github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	routeV1 := a.Group("/api/v1")

	routeV1.Post("/user/signup", controllers.SignUp)
	routeV1.Post("/user/signin", controllers.SignIn)
}
