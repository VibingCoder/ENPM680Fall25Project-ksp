package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"retrobytes/internal/repos"
	"retrobytes/internal/services"
	"retrobytes/internal/validate"
)

type OrderHandler struct {
	Cart  *services.CartService
	Order *services.OrderService
	Repo  *repos.OrderRepo
}

type OrderDeps struct {
	Cart *services.CartService
	Ord  *services.OrderService
}

func (h *OrderHandler) ensureSID(c *fiber.Ctx) string {
	sid := c.Cookies("sid")
	if sid == "" {
		sid = uuid.NewString()
		c.Cookie(&fiber.Cookie{
			Name:     "sid",
			Value:    sid,
			Path:     "/",
			HTTPOnly: true,
		})
	}
	return sid
}

func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	cv, err := h.Cart.View(h.ensureSID(c))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Render("checkout", fiber.Map{"Cart": cv})
}

func (h *OrderHandler) Place(c *fiber.Ctx) error {
	sid := h.ensureSID(c)

	// Validate region/ZIP
	region, ok := validate.Region(c.FormValue("region"))
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("invalid region/ZIP")
	}

	// Validate email and name
	email, ok := validate.Email(c.FormValue("email"))
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("invalid email")
	}
	name := strings.TrimSpace(c.FormValue("name"))
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("name is required")
	}

	// Normalize fulfillment
	fulfillment := strings.ToLower(strings.TrimSpace(c.FormValue("fulfillment")))
	if fulfillment != "delivery" && fulfillment != "pickup" {
		fulfillment = "delivery"
	}

	contact := services.Contact{Name: name, Email: email}

	orderID, err := h.Order.Place(sid, region, fulfillment, contact)
	if err != nil {
		// business rule errors (e.g., insufficient stock) surface as 400
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Show detailed confirmation page
	return c.Redirect("/order/" + orderID)
}

func (h *OrderHandler) View(c *fiber.Ctx) error {
	oid := c.Params("id")
	if oid == "" {
		return c.Status(fiber.StatusNotFound).Render("notfound", fiber.Map{"Message": "Order not found"})
	}

	o, items, err := h.Repo.Get(oid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("notfound", fiber.Map{"Message": "Order not found"})
	}

	return c.Render("order", fiber.Map{"Order": o, "Items": items})
}
