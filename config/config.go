package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Server     ServerConfig
	Database   DatabaseConfig
	Logging    LoggingConfig
	Mattermost MattermostConfig
}

type AppConfig struct {
	Name                string        `mapstructure:"name"`
	Version             string        `mapstructure:"version"`
	DefaultVoteDuration time.Duration `mapstructure:"default_vote_duration"`
	MaxPollOptions      int           `mapstructure:"max_poll_options"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Tarantool TarantoolConfig `mapstructure:"tarantool"`
}

type TarantoolConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	Timeout        time.Duration `mapstructure:"timeout"`
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
	MaxReconnects  int           `mapstructure:"max_reconnects"`
	Spaces         struct {
		Polls string `mapstructure:"polls"`
		Votes string `mapstructure:"votes"`
	} `mapstructure:"spaces"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MattermostConfig struct {
	BotToken    string `mapstructure:"bot_token"`
	BotID       string `mapstructure:"bot_id"`
	WebhookPath string `mapstructure:"webhook_path"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
