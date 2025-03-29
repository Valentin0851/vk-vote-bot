package domain

import "context"

type Repository interface {
	PollRepository
	VoteRepository
}

type PollRepository interface {
	Create(ctx context.Context, poll *Poll) error
	GetByID(ctx context.Context, id string) (*Poll, error)
	UpdateStatus(ctx context.Context, id string, status PollStatus) error
	Delete(ctx context.Context, id string) error
}

type VoteRepository interface {
	Create(ctx context.Context, vote *Vote) error
	GetByPollID(ctx context.Context, pollID string) ([]*Vote, error)
	DeleteByPollID(ctx context.Context, pollID string) error
}
