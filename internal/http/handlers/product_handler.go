package handlers

import (
	"retrobytes/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	Catalog *services.CatalogService
}

func (h *ProductHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := h.Catalog.GetProduct(id)
	if err != nil || p.ID == "" {
		return c.Status(404).Render("notfound", fiber.Map{"Message": "This item is no longer available"})
	}
	return c.Render("product", fiber.Map{"P": p})
}
