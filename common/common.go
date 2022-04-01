package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetenvData(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	value := os.Getenv(key)
	return value
}
