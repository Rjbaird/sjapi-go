package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT        string
	MONGODB_URI string
}

func ENV() (*Config, error) {
	godotenv.Load(".env")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		fmt.Println("no PORT environment variable provided")
		fmt.Println("Setting PORT to 3000")
		PORT = "3000"
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	config := Config{PORT: PORT, MONGODB_URI: MONGODB_URI}

	return &config, nil
}
