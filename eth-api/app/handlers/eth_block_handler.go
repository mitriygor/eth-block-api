package handlers

import (
	"eth-api/app/services"
	"github.com/gofiber/fiber/v2"
	"log"
)

type EthBlockHandler struct {
	EthBlockService services.EthBlockService
}

func NewEthBlockHandler(service services.EthBlockService) *EthBlockHandler {
	return &EthBlockHandler{
		EthBlockService: service,
	}
}

func (h *EthBlockHandler) GetBlockByIdentifierHandler(c *fiber.Ctx) error {
	blockIdentifier := c.Params("blockIdentifier")
	log.Printf("API::GetBlockByIdentifierHandler: %v", blockIdentifier)

	ethBlock, err := h.EthBlockService.GetBlockByIdentifierService(blockIdentifier)
	log.Printf("API::GetBlockByIdentifierHandler::ethBlock: %v; err: ; %v;", ethBlock, err)

	if err != nil {
		log.Printf("API::ERROR::GetBlockByIdentifierHandler::err: ; %v;", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}
