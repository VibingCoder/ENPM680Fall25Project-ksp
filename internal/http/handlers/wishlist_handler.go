package handlers

import (
	"retrobytes/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	Wish *services.WishlistService
}

func (h *WishlistHandler) ensureSID(c *fiber.Ctx) string {
	sid := c.Cookies("sid")
	if sid == "" {
		sid = uuid.NewString()
		c.Cookie(&fiber.Cookie{Name: "sid", Value: sid, Path: "/", HTTPOnly: true})
	}
	return sid
}

func (h *WishlistHandler) List(c *fiber.Ctx) error {
	sid := h.ensureSID(c)
	items, err := h.Wish.List(sid)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("wishlist", fiber.Map{"Items": items})
}

func (h *WishlistHandler) Save(c *fiber.Ctx) error {
	sid := h.ensureSID(c)
	pid := c.FormValue("productId")
	if pid == "" {
		return c.Status(400).SendString("missing productId")
	}
	if err := h.Wish.Save(sid, pid); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	// redirect back to product or wishlist
	back := c.Get("Referer")
	if back == "" {
		back = "/wishlist"
	}
	return c.Redirect(back)
}

func (h *WishlistHandler) Unsave(c *fiber.Ctx) error {
	sid := h.ensureSID(c)
	pid := c.FormValue("productId")
	if pid == "" {
		return c.Status(400).SendString("missing productId")
	}
	if err := h.Wish.Unsave(sid, pid); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Redirect("/wishlist")
}
