package tarantool

import (
	"errors"
	"fmt"
	"mattermost-vote-bot/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
)

const (
	pollsSpace       = "polls"
	votesSpace       = "votes"
	resultsSpace     = "results"
	pollsByCreator   = "polls_by_creator"
	activePollsIndex = "active_polls"
)

type TarantoolStorage struct {
	conn *tarantool.Connection
}

func NewStorage(address, user, password string) (*TarantoolStorage, error) {
	opts := tarantool.Opts{
		User:          user,
		Pass:          password,
		Timeout:       5 * time.Second,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
	}

	conn, err := tarantool.Connect(address, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tarantool: %w", err)
	}

	if _, err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("tarantool ping failed: %w", err)
	}

	return &TarantoolStorage{conn: conn}, nil
}

func (s *TarantoolStorage) CreatePoll(poll domain.Poll) (string, error) {
	if poll.ID == "" {
		poll.ID = uuid.New().String()
	}

	_, err := s.conn.Insert(pollsSpace, []interface{}{
		poll.ID,
		poll.CreatorID,
		poll.Question,
		poll.Options,
		poll.CreatedAt,
		true,
	})

	if err != nil {
		return "", fmt.Errorf("failed to insert poll: %w", err)
	}

	_, err = s.conn.Insert(pollsByCreator, []interface{}{
		poll.CreatorID,
		poll.ID,
	})

	if err != nil {
		s.conn.Delete(pollsSpace, []interface{}{poll.ID})
		return "", fmt.Errorf("failed to update creator index: %w", err)
	}

	return poll.ID, nil
}

func (s *TarantoolStorage) Vote(pollID, userID, option string) error {
	resp, err := s.conn.Select(pollsSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil {
		return fmt.Errorf("failed to select poll: %w", err)
	}
	if len(resp.Data) == 0 {
		return errors.New("poll does not exist")
	}

	pollData := resp.Data[0].([]interface{})
	isActive := pollData[5].(bool)
	if !isActive {
		return errors.New("poll is not active")
	}

	options := pollData[3].(map[interface{}]interface{})
	if _, ok := options[option]; !ok {
		return errors.New("invalid option")
	}

	resp, err = s.conn.Select(votesSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{pollID, userID})
	if err != nil {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}
	if len(resp.Data) > 0 {
		return errors.New("user already voted")
	}

	_, err = s.conn.Insert(votesSpace, []interface{}{
		pollID,
		userID,
		option,
		time.Now().Unix(),
	})

	if err != nil {
		return fmt.Errorf("failed to insert vote: %w", err)
	}

	err = s.updateResults(pollID, option)
	if err != nil {
		return fmt.Errorf("failed to update results: %w", err)
	}

	return nil
}

func (s *TarantoolStorage) updateResults(pollID, option string) error {
	_, err := s.conn.Call("box.atomic.counter_inc", []interface{}{
		resultsSpace,
		pollID,
		option,
		1,
	})

	return err
}

func (s *TarantoolStorage) GetResults(pollID string) (map[string]int, error) {
	resp, err := s.conn.Select(resultsSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil {
		return nil, fmt.Errorf("failed to select results: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no results found")
	}

	resultData := resp.Data[0].([]interface{})
	results := make(map[string]int)

	for k, v := range resultData[1].(map[interface{}]interface{}) {
		results[k.(string)] = int(v.(int64))
	}

	return results, nil
}

func (s *TarantoolStorage) EndPoll(pollID, userID string) error {
	resp, err := s.conn.Select(pollsSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil {
		return fmt.Errorf("failed to select poll: %w", err)
	}
	if len(resp.Data) == 0 {
		return errors.New("poll does not exist")
	}

	pollData := resp.Data[0].([]interface{})
	creatorID := pollData[1].(string)
	if creatorID != userID {
		return errors.New("only poll creator can end the poll")
	}

	_, err = s.conn.Update(pollsSpace, "primary", []interface{}{pollID}, []interface{}{
		[]interface{}{"=", 5, false},
	})

	return err
}

func (s *TarantoolStorage) DeletePoll(pollID, userID string) error {
	resp, err := s.conn.Select(pollsSpace, "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil {
		return fmt.Errorf("failed to select poll: %w", err)
	}
	if len(resp.Data) == 0 {
		return errors.New("poll does not exist")
	}

	pollData := resp.Data[0].([]interface{})
	creatorID := pollData[1].(string)
	if creatorID != userID {
		return errors.New("only poll creator can delete the poll")
	}

	if _, err := s.conn.Delete(pollsSpace, []interface{}{pollID}); err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	if _, err := s.conn.Delete(pollsByCreator, []interface{}{creatorID, pollID}); err != nil {
		return fmt.Errorf("failed to delete from creator index: %w", err)
	}

	if _, err := s.conn.Delete(resultsSpace, []interface{}{pollID}); err != nil {
		return fmt.Errorf("failed to delete results: %w", err)
	}

	return nil
}

func (s *TarantoolStorage) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
