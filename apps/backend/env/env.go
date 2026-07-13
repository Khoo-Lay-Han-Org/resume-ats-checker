package env

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		print("ENV file not found: " + err.Error())
	}
	return os.Getenv(key)
}
