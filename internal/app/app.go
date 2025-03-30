package app

import (
	"mattermost-vote-bot/internal/bot"
	"mattermost-vote-bot/internal/config"
	"mattermost-vote-bot/internal/infrastructure/mattermost"
	"mattermost-vote-bot/internal/infrastructure/storage/tarantool"
)

type App struct {
	bot *bot.Bot
}

func NewApp(cfg *config.Config) (*App, error) {
	mmClient, err := mattermost.NewClient(
		cfg.Mattermost.URL,
		cfg.Mattermost.Token,
		cfg.Mattermost.Team,
		cfg.Mattermost.Channel,
	)
	if err != nil {
		return nil, err
	}

	storage, err := tarantool.NewStorage(
		cfg.Tarantool.Address,
		cfg.Tarantool.User,
		cfg.Tarantool.Password,
	)
	if err != nil {
		return nil, err
	}

	return &App{
		bot: bot.NewBot(mmClient, storage),
	}, nil
}

func (a *App) Run() error {
	return a.bot.Start()
}
