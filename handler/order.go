package handler

import (
	"context"
	"log"
	"time"

	"github.com/b4cktr4ck5r3/micro-order/database"
	"github.com/b4cktr4ck5r3/micro-order/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllOrders(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var orders []model.Order

	cursor, err := orderCollection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Orders Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var order model.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": orders,
	})
}

func GetOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var order model.Order
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&order)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    order,
		"success": true,
	})
}

func AddOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new([]model.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	doc := make([]interface{}, len(*order))
	for i := 0; i < len(*order); i++ {
		doc[i] = (*order)[i]
	}

	result, err := orderCollection.InsertMany(ctx, doc)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Order inserted successfully",
	})

}

func UpdateOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(model.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": order,
	}
	_, err = orderCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Order updated successfully",
	})
}

func DeleteOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order not found",
			"error":   err,
		})
	}
	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Order deleted successfully",
	})
}
