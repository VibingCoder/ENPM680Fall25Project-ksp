package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"retrobytes/internal/services"
	"retrobytes/internal/validate"
)

type InventoryHandler struct {
	Inv *services.InventoryService
}

func (h *InventoryHandler) Check(c *fiber.Ctx) error {
	// Validate productId
	productID := strings.TrimSpace(c.Query("productId"))
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing productId",
		})
	}

	// Validate region/ZIP (allows simple ZIP/postal formats)
	region, ok := validate.Region(c.Query("region"))
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "enter a valid region/ZIP",
		})
	}

	// Business logic
	avail, err := h.Inv.CheckAvailability(productID, region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(avail)
}
