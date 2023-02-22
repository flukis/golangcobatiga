package main

import (
	"os"
	"pos-services/pkg/configs"
	"pos-services/pkg/routes"
	"pos-services/pkg/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// routes for public api
	routes.PublicRoutes(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
