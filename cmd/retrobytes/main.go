package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	html "github.com/gofiber/template/html/v2"

	"retrobytes/internal/config"
	"retrobytes/internal/http/handlers"
	"retrobytes/internal/repos"
)

func main() {
	cfg := config.Load()

	db, err := repos.OpenDB(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Templates
	engine := html.New("./web/templates", ".html")
	engine.Reload(true) // dev convenience; remove in prod
	app := fiber.New(fiber.Config{Views: engine})

	// Middleware (optional but helpful)
	app.Use(logger.New())
	app.Use(helmet.New())

	// Static
	app.Static("/static", "./web/static")
	app.Static("/media", cfg.MediaDir)

	// Deps
	deps := handlers.NewDeps(db, cfg)

	// Pages
	app.Get("/", deps.CategoryHandler.Home)
	app.Get("/search", deps.SearchHandler.Search)
	app.Get("/category/:id", deps.CategoryHandler.List)

	// Product routes
	app.Get("/product", func(c *fiber.Ctx) error {
		return c.Status(404).Render("notfound", fiber.Map{"Message": "This item is no longer available"})
	})
	app.Get("/product/:id", deps.ProductHandler.Detail)

	// API
	api := app.Group("/api/v1")
	api.Get("/availability", deps.InventoryHandler.Check)

	// Cart & Orders
	app.Get("/cart", deps.CartHandler.View)
	app.Post("/cart", deps.CartHandler.Add)
	app.Get("/checkout", deps.OrderHandler.Checkout)
	app.Post("/orders", deps.OrderHandler.Place)
	app.Get("/order/:id", deps.OrderHandler.View)

	// Wishlist
	app.Get("/wishlist", deps.WishlistHandler.List)
	app.Post("/wishlist", deps.WishlistHandler.Save)
	app.Post("/wishlist/delete", deps.WishlistHandler.Unsave)

	// Health
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) })

	// Catch-all 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("notfound", fiber.Map{
			"Message": "Page not found",
		})
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
