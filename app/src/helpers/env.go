package helpers

import (
	"cloudview/app/src/api/middleware/logger"
	"os"

	dotenv "github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	// load .env file
	if err := dotenv.Load(".env"); err != nil {
		logger.Logger.Error("Error loading .env file")
	}

	return os.Getenv(key)
}
