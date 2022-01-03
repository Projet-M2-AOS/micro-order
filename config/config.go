package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

//Retrieve env variable from .env file or system variable (in Docker) by key
func Config(key string) string {
	var err error = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return os.Getenv(key)
}
