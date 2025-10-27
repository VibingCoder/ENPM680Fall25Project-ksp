package handlers

import (
	"retrobytes/internal/services"
	"retrobytes/internal/validate"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CartHandler struct {
	Cart *services.CartService
}

func (h *CartHandler) ensureSID(c *fiber.Ctx) string {
	sid := c.Cookies("sid")
	if sid == "" {
		sid = uuid.NewString()
		c.Cookie(&fiber.Cookie{Name: "sid", Value: sid, Path: "/", HTTPOnly: true})
	}
	return sid
}

func (h *CartHandler) Add(c *fiber.Ctx) error {
	sid := h.ensureSID(c)
	productID := c.FormValue("productId")
	qty := validate.Qty(c.FormValue("qty"))

	//qty, _ := strconv.Atoi(c.FormValue("qty"))

	if qty <= 0 {
		qty = 1
	}
	if productID == "" {
		return c.Status(400).SendString("missing productId")
	}
	if err := h.Cart.Add(sid, productID, qty); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Redirect("/cart")
}

func (h *CartHandler) View(c *fiber.Ctx) error {
	sid := h.ensureSID(c)
	cv, err := h.Cart.View(sid)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("cart", fiber.Map{"Cart": cv})

}
