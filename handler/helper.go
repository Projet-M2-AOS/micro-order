package handler

import "github.com/gofiber/fiber/v2"

//Create error message for response with format like Nest.js standard
func makeErrorMsg(statusCode int, message string, errorMsg string) fiber.Map {
	return fiber.Map{
		"statusCode": statusCode,
		"message":    message,
		"error":      errorMsg,
	}
}
