package tarantool

import (
	"context"
	"fmt"
	"time"

	"mattermost-vote-bot/internal/domain"

	"github.com/viciious/go-tarantool"
)

type Repository struct {
	conn   *tarantool.Connection
	config *Config
}

type Config struct {
	Host           string
	Port           int
	User           string
	Password       string
	Timeout        time.Duration
	ReconnectDelay time.Duration
	MaxReconnects  int
	Spaces         struct {
		Polls string
		Votes string
	}
}

func New(cfg *Config) (*Repository, error) {
	conn, err := tarantool.Connect(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		tarantool.Options{
			User:           cfg.User,
			Password:       cfg.Password,
			Timeout:        cfg.Timeout,
			ReconnectDelay: cfg.ReconnectDelay,
			MaxReconnects:  cfg.MaxReconnects,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tarantool: %w", err)
	}

	return &Repository{
		conn:   conn,
		config: cfg,
	}, nil
}

func (r *Repository) Create(ctx context.Context, poll *domain.Poll) error {
	_, err := r.conn.Insert(r.config.Spaces.Polls, []interface{}{
		poll.ID,
		poll.Question,
		poll.Creator,
		poll.ChannelID,
		poll.CreatedAt.Unix(),
		string(poll.Status),
		poll.Options,
	})
	if err != nil {
		return fmt.Errorf("failed to insert poll: %w", err)
	}
	return nil
}

func (r *Repository) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
