package models

import (
	"errors"
	"time"
)

// Vote - доменная модель голоса
type Vote struct {
	PollID   string    `json:"poll_id"`   // ID опроса
	UserID   string    `json:"user_id"`   // ID пользователя
	OptionID string    `json:"option_id"` // ID выбранного варианта
	VotedAt  time.Time `json:"voted_at"`  // Время голосования
}

// Validate - валидация данных голоса
func (v *Vote) Validate() error {
	if v.PollID == "" {
		return ErrEmptyPollID
	}
	if v.UserID == "" {
		return ErrEmptyUserID
	}
	if v.OptionID == "" {
		return ErrEmptyOptionID
	}
	return nil
}

// Ошибки валидации
var (
	ErrEmptyPollID   = errors.New("poll ID cannot be empty")
	ErrEmptyUserID   = errors.New("user ID cannot be empty")
	ErrEmptyOptionID = errors.New("option ID cannot be empty")
)
