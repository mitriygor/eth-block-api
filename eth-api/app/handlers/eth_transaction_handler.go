package handlers

import (
	"eth-api/app/services"
	"github.com/gofiber/fiber/v2"
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
	ethTransaction, err := h.EthTransactionService.GetTransactionByHashService(hash)

	if err != nil {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthEvents",
		})
	}

	return c.JSON(events)
}
