package routes

import (
	"eth-api/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, ethBlockHandler *handlers.EthBlockHandler, ethTransactionHandler *handlers.EthTransactionHandler) {
	app.Get("/eth-blocks/:blockIdentifier", ethBlockHandler.GetBlockByIdentifierHandler)
	app.Get("/eth-transactions/:transactionHash", ethTransactionHandler.GetTransactionByHashHandler)
	app.Get("/eth-events/:address", ethTransactionHandler.GetTransactionsByAddressHandler)
}
