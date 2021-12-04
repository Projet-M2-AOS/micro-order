package router

import (
	"github.com/b4cktr4ck5r3/micro-order/handler"
	"github.com/gofiber/fiber/v2"
)

func OrderRoute(route fiber.Router) {
	route.Get("/", handler.GetAllOrders)
	route.Get("/:id", handler.GetOrder)
	route.Post("/", handler.AddOrder)
	route.Put("/:id", handler.UpdateOrder)
	route.Delete("/:id", handler.DeleteOrder)
}
