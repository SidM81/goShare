package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SidM81/goShare/config"
	"github.com/SidM81/goShare/routes"
	"github.com/joho/godotenv"
)

func main() {
	//load Configurations
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.MustHaveConfig()
	config.InitMinio(cfg)
	config.InitDatabase()
	// database Setup

	r := routes.SetupRouter()

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
