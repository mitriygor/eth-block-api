package handlers

import (
	"eth-api/app/helpers/logger"
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

func (h *EthBlockHandler) GetLatestEthBlocksHandler(c *fiber.Ctx) error {
	ethBlock, err := h.EthBlockService.GetLatestEthBlocks()

	if err != nil {
		logger.Error("eth-api::ERROR::GetBlockByIdentifierHandler", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}

func (h *EthBlockHandler) GetBlockByIdentifierHandler(c *fiber.Ctx) error {
	blockIdentifier := c.Params("blockIdentifier")

	ethBlock, err := h.EthBlockService.GetBlockByIdentifierService(blockIdentifier)

	if err != nil {
		logger.Error("eth-api::ERROR::GetBlockByIdentifierHandler", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}
