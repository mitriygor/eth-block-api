package routes

import (
	"eth-api/app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *handlers.EthBlockHandler) {
	app.Get("/eth-blocks/:blockNumber", handler.GetEthBlockByNumberHandler)
	app.Post("/eth-blocks", handler.CreateEthBlockHandler)
}

//eth_blockNumber — Returns the number of most recent block.
//eth_getBlockByNumber — Returns information about a block by block number.
//eth_getTransactionByHash — Returns the information about a transaction requested by transaction hash.
//eth_getLogs — Returns an array of all logs matching a given filter object.
