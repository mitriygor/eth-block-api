package handlers

import (
	"eth-api/app/services"
	"fmt"
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

func (h *EthBlockHandler) GetLatestEthBlocksHandler(c *fiber.Ctx) error {

	fmt.Printf("eth-api::EthBlockHandler::GetLatestEthBlocksHandler")

	ethBlock, err := h.EthBlockService.GetLatestEthBlocks()

	if err != nil {
		log.Printf("eth-api::ERROR::GetLatestBlocksHandler::err: ; %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}

func (h *EthBlockHandler) GetBlockByIdentifierHandler(c *fiber.Ctx) error {
	blockIdentifier := c.Params("blockIdentifier")
	log.Printf("eth-api::GetBlockByIdentifierHandler: %v\n", blockIdentifier)

	ethBlock, err := h.EthBlockService.GetBlockByIdentifierService(blockIdentifier)
	log.Printf("eth-api::GetBlockByIdentifierHandler::ethBlock: %v; err: ; %v\n", ethBlock, err)

	if err != nil {
		log.Printf("eth-api::ERROR::GetBlockByIdentifierHandler::err: ; %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving EthBlock",
		})
	}

	return c.JSON(ethBlock)
}
