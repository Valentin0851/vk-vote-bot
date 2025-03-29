package main

import (
	"log"
	"mattermost-vote-bot/config"
	"mattermost-vote-bot/internal/api"
	"mattermost-vote-bot/internal/infra/tarantool"
	"mattermost-vote-bot/internal/service"
	"mattermost-vote-bot/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)

	tarantoolRepo, err := tarantool.New(&tarantool.Config{
		Host:           cfg.Database.Tarantool.Host,
		Port:           cfg.Database.Tarantool.Port,
		User:           cfg.Database.Tarantool.User,
		Password:       cfg.Database.Tarantool.Password,
		Timeout:        cfg.Database.Tarantool.Timeout,
		ReconnectDelay: cfg.Database.Tarantool.ReconnectDelay,
		MaxReconnects:  cfg.Database.Tarantool.MaxReconnects,
		Spaces: struct {
			Polls string
			Votes string
		}{
			Polls: cfg.Database.Tarantool.Spaces.Polls,
			Votes: cfg.Database.Tarantool.Spaces.Votes,
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize Tarantool: %v", err)
	}
	defer tarantoolRepo.Close()

	pollService := service.NewPollService(
		tarantoolRepo,
		service.NewUUIDGenerator(),
		cfg.Database.Tarantool.Timeout,
	)

	server := api.New(
		&api.Config{
			Host:            cfg.Server.Host,
			Port:            cfg.Server.Port,
			ReadTimeout:     cfg.Server.ReadTimeout,
			WriteTimeout:    cfg.Server.WriteTimeout,
			ShutdownTimeout: cfg.Server.ShutdownTimeout,
		},
		pollService,
		log,
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Infof("Server started on %s:%d", cfg.Server.Host, cfg.Server.Port)
	<-done
	log.Info("Server stopped")
}
