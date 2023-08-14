package routes

import (
	"eth-api/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *handlers.EthBlockHandler) {
	app.Get("/eth-blocks/:blockIdentifier", handler.GetBlockByIdentifierHandler)
}
