package handlers

import (
	"strings"

	"retrobytes/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	Catalog *services.CatalogService
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	q := strings.ToLower(strings.TrimSpace(c.Query("q")))
	category := strings.TrimSpace(c.Query("category"))
	condition := strings.TrimSpace(c.Query("condition")) // FIRST_HAND | SECOND_HAND

	products, err := h.Catalog.Search(q, category, condition, 1, 20)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Render("search", fiber.Map{
		"Q": q, "CategoryID": category, "Condition": condition,
		"Products": products, "Count": len(products),
	})
}
