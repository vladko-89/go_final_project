package main

import (
	"final-project/internal/static"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	todoPort := os.Getenv("TODO_PORT")
	if todoPort == "" {
		fmt.Println("TODO_PORT is not set")
		todoPort = ":7540"
	}

	router := http.NewServeMux()

	static.NewStaticHandler(router)

	server := http.Server{
		Addr:    todoPort,
		Handler: router,
	}

	fmt.Println("Server is running on port 8081")
	server.ListenAndServe()
}
