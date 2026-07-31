package main

import (
	"fmt"
	"log"

	"github.com/akansha204/cli-auth/internal/auth"
	"github.com/akansha204/cli-auth/internal/cli"
	"github.com/akansha204/cli-auth/internal/config"
	"github.com/akansha204/cli-auth/internal/database"
	"github.com/akansha204/cli-auth/internal/repository"
	"github.com/akansha204/cli-auth/internal/session"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if err := database.Initialize(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(database.DB)
	sessionRepo := repository.NewSessionRepository(database.DB)

	authService := auth.NewAuthService(userRepo)
	sessionManager := session.NewManager(sessionRepo, config.AppConfig.SessionTimeout)

	app := cli.NewApp(authService, sessionManager)

	fmt.Println("Welcome to CLI Auth")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
