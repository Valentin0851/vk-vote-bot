package service

import (
	"context"
	"errors"
	"time"

	"mattermost-vote-bot/internal/domain"
	"mattermost-vote-bot/internal/domain/repository"
)

// VoteService - сервис для работы с голосованиями
type VoteService struct {
	repo           repository.VoteRepository
	pollRepo       repository.PollRepository
	defaultTimeout time.Duration
}

// NewVoteService - конструктор VoteService
func NewVoteService(
	repo repository.VoteRepository,
	pollRepo repository.PollRepository,
	timeout time.Duration,
) *VoteService {
	return &VoteService{
		repo:           repo,
		pollRepo:       pollRepo,
		defaultTimeout: timeout,
	}
}

// VoteRequest - DTO для голосования
type VoteRequest struct {
	PollID   string `json:"poll_id"`
	UserID   string `json:"user_id"`
	OptionID string `json:"option_id"`
}

// Vote - обработка голосования
func (s *VoteService) Vote(ctx context.Context, req *VoteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.defaultTimeout)
	defer cancel()

	// 1. Проверяем существование опроса
	poll, err := s.pollRepo.GetByID(ctx, req.PollID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPollNotFound
		}
		return err
	}

	// 2. Проверяем статус опроса
	if poll.Status != domain.PollStatusActive {
		return ErrPollNotActive
	}

	// 3. Проверяем существование варианта
	if _, exists := poll.Options[req.OptionID]; !exists {
		return ErrInvalidOption
	}

	// 4. Создаем голос
	vote := &domain.Vote{
		PollID:   req.PollID,
		UserID:   req.UserID,
		OptionID: req.OptionID,
		VotedAt:  time.Now(),
	}

	if err := vote.Validate(); err != nil {
		return err
	}

	// 5. Сохраняем в репозитории
	return s.repo.Create(ctx, vote)
}

// ResultsResponse - DTO для результатов голосования
type ResultsResponse struct {
	PollID     string            `json:"poll_id"`
	Question   string            `json:"question"`
	Options    map[string]string `json:"options"`
	VoteCounts map[string]int    `json:"vote_counts"`
	TotalVotes int               `json:"total_votes"`
	UserVote   string            `json:"user_vote,omitempty"`
}

// GetResults - получение результатов голосования
func (s *VoteService) GetResults(ctx context.Context, pollID, userID string) (*ResultsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.defaultTimeout)
	defer cancel()

	// 1. Получаем опрос
	poll, err := s.pollRepo.GetByID(ctx, pollID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrPollNotFound
		}
		return nil, err
	}

	// 2. Получаем все голоса
	votes, err := s.repo.GetByPollID(ctx, pollID)
	if err != nil {
		return nil, err
	}

	// 3. Подсчет результатов
	results := &ResultsResponse{
		PollID:     poll.ID,
		Question:   poll.Question,
		Options:    poll.Options,
		VoteCounts: make(map[string]int),
	}

	// 4. Считаем голоса
	for _, vote := range votes {
		results.VoteCounts[vote.OptionID]++
		if vote.UserID == userID {
			results.UserVote = vote.OptionID
		}
	}

	// 5. Считаем общее количество
	for _, count := range results.VoteCounts {
		results.TotalVotes += count
	}

	return results, nil
}

// Ошибки сервиса
var (
	ErrPollNotFound  = errors.New("poll not found")
	ErrPollNotActive = errors.New("poll is not active")
	ErrInvalidOption = errors.New("invalid option ID")
)
