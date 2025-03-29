package models

import (
	"errors"
	"time"
)

type Poll struct {
	ID        string
	Question  string
	Creator   string
	ChannelID string
	CreatedAt time.Time
	Status    PollStatus
	Options   map[string]string
}

type PollStatus string

const (
	PollStatusActive PollStatus = "active"
	PollStatusEnded  PollStatus = "ended"
)

var (
	ErrEmptyQuestion    = errors.New("poll question cannot be empty")
	ErrNotEnoughOptions = errors.New("poll must have at least 2 options")
)

func (p *Poll) Validate() error {
	if p.Question == "" {
		return ErrEmptyQuestion
	}
	if len(p.Options) < 2 {
		return ErrNotEnoughOptions
	}
	return nil
}
