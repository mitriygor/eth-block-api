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
	log.Printf("eth-api::INFO::GetTransactionByHashHandler::hash: %v\n", hash)

	ethTransaction, err := h.EthTransactionService.GetTransactionByHashService(hash)

	log.Printf("eth-api::INFO::GetTransactionByHashHandler::ethTransaction: %v\n", ethTransaction)
	log.Printf("eth-api::INFO::GetTransactionByHashHandler::err: %v\n", err)

	if err != nil {
		log.Printf("eth-api::ERROR::GetTransactionByHashHandler::err: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthTransaction",
		})
	}

	return c.JSON(ethTransaction)
}

func (h *EthTransactionHandler) GetTransactionsByAddressHandler(c *fiber.Ctx) error {

	log.Printf("eth-api::INFO::GetTransactionsByAddressHandler::address: %v\n", c.Params("address"))
	address := c.Params("address")

	events, err := h.EthTransactionService.GetTransactionsByAddressService(address)

	if err != nil {
		log.Printf("eth-api::ERROR::GetEventsByAddressHandler::err: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthEvents",
		})
	}

	return c.JSON(events)
}
