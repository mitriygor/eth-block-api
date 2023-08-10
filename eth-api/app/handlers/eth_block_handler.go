package handlers

import (
	"eth-api/app/models"
	"eth-api/app/services"
	"github.com/gofiber/fiber/v2"
)

type EthBlockHandler struct {
	EthBlockService services.EthBlockService
}

func NewEthBlockHandler(service services.EthBlockService) *EthBlockHandler {
	return &EthBlockHandler{
		EthBlockService: service,
	}
}

func (h *EthBlockHandler) GetEthBlockByNumberHandler(c *fiber.Ctx) error {
	blockNumber := c.Params("blockNumber")

	ethBlock, err := h.EthBlockService.GetEthBlockByNumberService(blockNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}

func (h *EthBlockHandler) CreateEthBlockHandler(c *fiber.Ctx) error {
	var createEthBlockDto models.CreateEthBlockDto
	err := c.BodyParser(&createEthBlockDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Save the new EthBlock to the repository
	ethBlock, err := h.EthBlockService.CreateEthBlockService(createEthBlockDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating EthBlock",
		})
	}

	return c.JSON(ethBlock)
}
