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
		return c.Status(fiber.StatusNotFound).JSON(
			makeErrorMsg(fiber.StatusNotFound, err, "Orders not found"))
	}

	for cursor.Next(ctx) {
		var order model.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

func GetOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var order model.Order
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err, "Order not found"))
	}

	err = findResult.Decode(&order)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err, "Error during decode order"))
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func AddOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new([]model.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err, "Failed to parse body"))
	}

	doc := make([]interface{}, len(*order))
	for i := 0; i < len(*order); i++ {
		doc[i] = (*order)[i]
	}

	result, err := orderCollection.InsertMany(ctx, doc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err, "Failed to insert order"))
	}

	var insertedOrder []model.Order
	for _, element := range result.InsertedIDs {
		var currentOrder model.Order
		findResult := orderCollection.FindOne(ctx, bson.M{"_id": element})
		if err := findResult.Err(); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err, "Inserted order not found"))
		}

		err = findResult.Decode(&currentOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err, "Error during decode order"))
		}

		insertedOrder = append(insertedOrder, currentOrder)
	}

	return c.Status(fiber.StatusCreated).JSON(insertedOrder)
}

func UpdateOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(model.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err, "Failed to parse body"))
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err, "Order not found"))
	}

	update := bson.M{
		"$set": order,
	}
	_, err = orderCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err, "Error during updating order"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func DeleteOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err, "Order not found"))
	}
	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err, "Error during deleting order"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
