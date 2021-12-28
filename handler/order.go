package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/b4cktr4ck5r3/micro-order/database"
	"github.com/b4cktr4ck5r3/micro-order/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Get all orders.
// @Description Return all orders.
// @Tags Orders
// @Produce json
// @Success 200 {array} model.Order
// @Router /orders [get]
func GetAllOrders(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var orders []model.Order

	cursor, err := orderCollection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			makeErrorMsg(fiber.StatusNotFound, err.Error(), "Orders not found"))
	}

	for cursor.Next(ctx) {
		var order model.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

// @Summary Get one orders.
// @Description Return one orders.
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Router /orders/{id} [get]
func GetOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var order model.Order
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	err = findResult.Decode(&order)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Error during decode order"))
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

// @Summary Create a new order
// @Description Create a new order with the input payload
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Order body model.Order true "Create order"
// @Success 200 {array} model.Order
// @Router /orders [post]
func AddOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new([]model.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Failed to parse body"))
	}

	doc := make([]interface{}, len(*order))
	validationErrors := ""
	for i := 0; i < len(*order); i++ {
		errors := model.ValidateStruct((*order)[i])
		if errors != "" {
			validationErrors += errors
		}
		doc[i] = (*order)[i]
	}

	if validationErrors != "" {
		return c.Status(fiber.StatusBadRequest).JSON((makeErrorMsg(fiber.StatusBadRequest, validationErrors, "Bad Request")))
	}

	if len(doc) == 0 {
		return c.Status(fiber.StatusCreated).JSON(doc)
	}

	result, err := orderCollection.InsertMany(ctx, doc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err.Error(), "Failed to insert order"))
	}

	var insertedOrder []model.Order
	for _, element := range result.InsertedIDs {
		var currentOrder model.Order
		findResult := orderCollection.FindOne(ctx, bson.M{"_id": element})
		if err := findResult.Err(); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Inserted order not found"))
		}

		err = findResult.Decode(&currentOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err.Error(), "Error during decode order"))
		}

		insertedOrder = append(insertedOrder, currentOrder)
	}

	return c.Status(fiber.StatusCreated).JSON(insertedOrder)
}

// @Summary Update order
// @Description Update order with the input payload
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param Order body model.Order true "Update order"
// @Param id path string true "Order ID"
// @Success 201
// @Router /orders/{id} [put]
func UpdateOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(model.Order)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Failed to parse body"))
	}

	update := bson.M{
		"$set": order,
	}

	result, err := orderCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Error during updating order"))
	}

	if &result.ModifiedCount == nil {
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Error during updating order"))
	}

	updatedResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	var updatedOrder model.Order
	err = updatedResult.Decode(&updatedOrder)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Error during decode order"))
	}

	fmt.Println(updatedOrder)
	return c.Status(fiber.StatusNoContent).JSON(updatedOrder)
}

// @Summary Delete order
// @Description Delete order
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 201
// @Router /orders/{id} [delete]
func DeleteOrder(c *fiber.Ctx) error {
	orderCollection := database.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "ID is not an Object ID"))
	}

	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err.Error(), "Error during deleting order"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
