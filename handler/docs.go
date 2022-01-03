package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
)

//Retrieve swagger json file created after using the command swag init (from swaggo module) as response
func GetSwaggerJson(c *fiber.Ctx) error {
	//Try to open generated swagger.json file
	jsonFile, err := os.Open("docs/swagger.json")
	if err != nil {
		fmt.Println(err)
		return c.Status(404).JSON(makeErrorMsg(404, "config not found", err.Error()))
	}

	defer jsonFile.Close()
	//Read file if opened
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var anyJson map[string]interface{}
	//Deserialize bytes of readed file ton anonymous structure (map[string]interface{})
	json.Unmarshal(byteValue, &anyJson)
	return c.Status(200).JSON(anyJson)
}
