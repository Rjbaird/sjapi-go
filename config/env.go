package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func ENV() (*Config, error) {
	godotenv.Load(".env")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		fmt.Println("no PORT environment variable provided")
		fmt.Println("Setting PORT to 3000")
		PORT = "3000"
	}

	config := Config{PORT: PORT}

	return &config, nil
}
