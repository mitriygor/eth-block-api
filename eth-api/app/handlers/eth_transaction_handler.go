package handlers

import (
	"eth-api/app/services"
	"github.com/gofiber/fiber/v2"
	"log"
)

type EthTransactionHandler struct {
	EthTransactionService services.EthTransactionService
}

func NewEthTransactionHandler(service services.EthTransactionService) *EthTransactionHandler {
	return &EthTransactionHandler{
		EthTransactionService: service,
	}
}

func (h *EthTransactionHandler) GetTransactionByHashHandler(c *fiber.Ctx) error {

	hash := c.Params("transactionHash")
	log.Printf("API::INFO::GetTransactionByHashHandler::hash: %v\n", hash)

	ethTransaction, err := h.EthTransactionService.GetTransactionByHashService(hash)

	log.Printf("API::INFO::GetTransactionByHashHandler::ethTransaction: %v\n", ethTransaction)
	log.Printf("API::INFO::GetTransactionByHashHandler::err: %v\n", err)

	if err != nil {
		log.Printf("API::ERROR::GetTransactionByHashHandler::err: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthTransaction",
		})
	}

	return c.JSON(ethTransaction)
}

func (h *EthTransactionHandler) GetTransactionsByAddressHandler(c *fiber.Ctx) error {
	address := c.Params("address")

	events, err := h.EthTransactionService.GetTransactionsByAddressService(address)

	if err != nil {
		log.Printf("API::ERROR::GetEventsByAddressHandler::err: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthEvents",
		})
	}

	return c.JSON(events)
}
