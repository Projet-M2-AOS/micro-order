package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetSwaggerJson(c *fiber.Ctx) error {
	jsonFile, err := os.Open("docs/swagger.json")
	if err != nil {
		fmt.Println(err)
		return c.Status(404).JSON(makeErrorMsg(404, "config not found", err.Error()))
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var anyJson map[string]interface{}
	json.Unmarshal(byteValue, &anyJson)
	return c.Status(200).JSON(anyJson)
}
