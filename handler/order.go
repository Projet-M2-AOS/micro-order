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

// @Summary Get all orders.
// @Description Return all orders.
// @Tags micro-orders
// @Param user query string false "search by userid"
// @Produce json
// @Success 200 {array} model.Order
// @Router /orders [get]
func GetAllOrders(c *fiber.Ctx) error {
	//Get mongo collection
	orderCollection := database.MI.DB.Collection("orders")
	//Set database context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var orders []model.Order
	//Try to retrieve userId from query parameter
	userId, err := primitive.ObjectIDFromHex(c.Query("userId"))

	filter := bson.M{}
	//Case when orders are query by user
	if err == nil {
		filter = bson.M{"user": bson.M{"$eq": userId}}
	}

	//Find documents in filtered collection
	cursor, err := orderCollection.Find(ctx, filter)
	defer cursor.Close(ctx)

	//Check if the find method didn't have error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			makeErrorMsg(fiber.StatusNotFound, err.Error(), "Orders not found"))
	}

	//Add all find order
	for cursor.Next(ctx) {
		var order model.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}

	//Case of empty collection (need to return empty array, like Nest.js standard : without this code Go return "nil" instead)
	if len(orders) == 0 {
		return c.Status(fiber.StatusOK).JSON([]model.Order{})
	}

	//Return orders
	return c.Status(fiber.StatusOK).JSON(orders)
}

// @Summary Get one orders.
// @Description Return one orders.
// @Tags micro-orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Router /orders/{id} [get]
func GetOrder(c *fiber.Ctx) error {
	//Get mongo collection
	orderCollection := database.MI.DB.Collection("orders")
	//Set database context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var order model.Order
	//Check if given ID is on ObjectID format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	//Find the order
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	//Check for error during finding order
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	//Try to decode to struct
	err = findResult.Decode(&order)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Error during decode order"))
	}

	//Return the order
	return c.Status(fiber.StatusOK).JSON(order)
}

// @Summary Create a new order
// @Description Create a new order with the input payload
// @Tags micro-orders
// @Accept  json
// @Produce  json
// @Param Order body model.Order true "Create order"
// @Success 200 {array} model.Order
// @Router /orders [post]
func AddOrder(c *fiber.Ctx) error {
	//Get mongo collection
	orderCollection := database.MI.DB.Collection("orders")
	//Set database context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new([]model.Order)

	//Try to parse Body to validate schema
	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Failed to parse body"))
	}

	//Validate each item of the body array and save validation error to a string
	doc := make([]interface{}, len(*order))
	validationErrors := ""
	for i := 0; i < len(*order); i++ {
		errors := model.ValidateStruct((*order)[i])
		if errors != "" {
			validationErrors += errors
		}
		doc[i] = (*order)[i]
	}

	//If something went wrong during validation, return error to user with explicit message
	if validationErrors != "" {
		return c.Status(fiber.StatusBadRequest).JSON((makeErrorMsg(fiber.StatusBadRequest, validationErrors, "Bad Request")))
	}

	//Case when the array is empty, return empty array like Nest.js
	if len(doc) == 0 {
		return c.Status(fiber.StatusCreated).JSON(doc)
	}

	//Insert to database
	result, err := orderCollection.InsertMany(ctx, doc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err.Error(), "Failed to insert order"))
	}

	//Retrieve inserted element to return them in the response
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
// @Tags micro-orders
// @Accept  json
// @Produce  json
// @Param Order body model.Order true "Update order"
// @Param id path string true "Order ID"
// @Success 201
// @Router /orders/{id} [put]
func UpdateOrder(c *fiber.Ctx) error {
	//Get mongo collection
	orderCollection := database.MI.DB.Collection("orders")
	//Set database context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(model.Order)

	//Check if given ID is on ObjectID format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	//Try to find the order to update
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	//Try to parse body
	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Failed to parse body"))
	}

	//Mongo update
	update := bson.M{
		"$set": order,
	}

	//Try to update the document
	result, err := orderCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Error during updating order"))
	}

	//Check if document has been updated
	if &result.ModifiedCount == nil {
		return c.Status(fiber.StatusBadRequest).JSON(makeErrorMsg(fiber.StatusBadRequest, err.Error(), "Error during updating order"))
	}

	//Retrieve updated document to return it in the response
	updatedResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Order not found"))
	}

	var updatedOrder model.Order
	err = updatedResult.Decode(&updatedOrder)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "Error during decode order"))
	}

	return c.Status(fiber.StatusNoContent).JSON(updatedOrder)
}

// @Summary Delete order
// @Description Delete order
// @Tags micro-orders
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 201
// @Router /orders/{id} [delete]
func DeleteOrder(c *fiber.Ctx) error {
	//Get mongo collection
	orderCollection := database.MI.DB.Collection("orders")
	//Get mongo collection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Check if given ID is on ObjectID format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(makeErrorMsg(fiber.StatusNotFound, err.Error(), "ID is not an Object ID"))
	}

	//Try to delete document in collection
	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(makeErrorMsg(fiber.StatusInternalServerError, err.Error(), "Error during deleting order"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
