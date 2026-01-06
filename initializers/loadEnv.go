package initializers

import (
	"log"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Println("Warning: .env.local file not found, using system environment variables")
	}
}
