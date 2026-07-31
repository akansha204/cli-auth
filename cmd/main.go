package main

import (
	"fmt"
	"log"

	"github.com/akansha204/cli-auth/internal/config"
	"github.com/akansha204/cli-auth/internal/database"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if err := database.Initialize(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Welcome to CLI Auth")
}
