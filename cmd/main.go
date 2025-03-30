package main

import (
	"Mattermost-bot-VK/internal/app"
	"Mattermost-bot-VK/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
