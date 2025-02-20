package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	Port PortCongig
}

type DbConfig struct {
	Path string
}

type PortCongig struct {
	Port string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("TODO_PORT")

	if port == "" {
		fmt.Println("TODO_PORT is not set")
		port = "7540"
	}

	return &Config{
		Db: DbConfig{
			Path: os.Getenv("TODO_DBFILE"),
		},
		Port: PortCongig{
			Port: fmt.Sprintf(":%s", port),
		},
	}
}
