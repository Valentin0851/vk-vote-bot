package service

import (
	"context"
	"time"

	"mattermost-vote-bot/internal/domain"
)

type PollService struct {
	repo    domain.PollRepository
	idGen   IDGenerator
	timeout time.Duration
}

func NewPollService(repo domain.PollRepository, idGen IDGenerator, timeout time.Duration) *PollService {
	return &PollService{
		repo:    repo,
		idGen:   idGen,
		timeout: timeout,
	}
}

func (s *PollService) Create(ctx context.Context, req *CreatePollRequest) (*domain.Poll, error) {
	poll := &domain.Poll{
		ID:        s.idGen.Generate(),
		Question:  req.Question,
		Creator:   req.Creator,
		ChannelID: req.ChannelID,
		CreatedAt: time.Now(),
		Status:    domain.PollStatusActive,
		Options:   req.Options,
	}

	if err := poll.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.repo.Create(ctx, poll); err != nil {
		return nil, err
	}

	return poll, nil
}
