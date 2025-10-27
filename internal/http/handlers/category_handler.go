package handlers

import (
	"retrobytes/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	Catalog *services.CatalogService
}

func (h *CategoryHandler) Home(c *fiber.Ctx) error {
	cats, err := h.Catalog.ListCategories()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("home", fiber.Map{"Categories": cats})

}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	catID := c.Params("id")
	products, err := h.Catalog.ListProductsByCategory(catID, 1, 12)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("category", fiber.Map{"CategoryID": catID, "Products": products})

}
