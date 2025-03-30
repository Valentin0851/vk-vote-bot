package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"Mattermost-bot-VK/internal/domain"

	"github.com/google/uuid"
)

type PollService struct {
	repo domain.PollRepository
	bot  *Bot
}

func NewPollService(repo domain.PollRepository) *PollService {
	return &PollService{repo: repo}
}

func (ps *PollService) SetBot(bot *Bot) {
	ps.bot = bot
}

func (ps *PollService) handleCreatePoll(userID, channelID, args string) error {
	parts := strings.SplitN(args, "\"", -1)
	if len(parts) < 5 {
		return ps.bot.SendMessage("Invalid format. Usage: /poll create \"Question\" \"Option1\" \"Option2\" ...")
	}

	question := strings.TrimSpace(parts[1])
	options := make(map[string]string)
	for i, opt := range parts[3:] {
		if i%2 == 0 {
			opt = strings.TrimSpace(opt)
			if opt != "" {
				options[fmt.Sprintf("opt%d", i/2+1)] = opt
			}
		}
	}

	if len(options) < 2 {
		return ps.bot.SendMessage("Poll must have at least 2 options")
	}

	poll := domain.Poll{
		ID:        uuid.New().String(),
		CreatorID: userID,
		Question:  question,
		Options:   options,
		CreatedAt: time.Now().Unix(),
		IsActive:  true,
	}

	pollID, err := ps.repo.CreatePoll(poll)
	if err != nil {
		log.Printf("Failed to create poll: %v", err)
		return ps.bot.SendMessage("Failed to create poll")
	}

	optionsText := ""
	for id, text := range options {
		optionsText += fmt.Sprintf("%s: %s\n", id, text)
	}

	return ps.bot.SendMessage(fmt.Sprintf(
		"Poll created successfully!\n"+
			"Question: %s\n"+
			"Options:\n%s\n"+
			"To vote: /poll vote %s <optionId>",
		question, optionsText, pollID))
}

func (ps *PollService) handleVote(userID, args string) error {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		return ps.bot.SendMessage("Invalid format. Usage: /poll vote <pollId> <optionId>")
	}

	pollID := parts[0]
	option := parts[1]

	if err := ps.repo.Vote(pollID, userID, option); err != nil {
		if errors.Is(err, domain.ErrPollNotFound) {
			return ps.bot.SendMessage("Poll not found")
		}
		if errors.Is(err, domain.ErrPollNotActive) {
			return ps.bot.SendMessage("This poll is no longer active")
		}
		if errors.Is(err, domain.ErrInvalidOption) {
			return ps.bot.SendMessage("Invalid option selected")
		}
		if errors.Is(err, domain.ErrAlreadyVoted) {
			return ps.bot.SendMessage("You have already voted in this poll")
		}
		log.Printf("Vote error: %v", err)
		return ps.bot.SendMessage("Failed to process your vote")
	}

	return ps.bot.SendMessage("Your vote has been recorded!")
}

func (ps *PollService) handleResults(args string) error {
	pollID := strings.TrimSpace(args)
	if pollID == "" {
		return ps.bot.SendMessage("Please specify poll ID")
	}

	results, err := ps.repo.GetResults(pollID)
	if err != nil {
		if errors.Is(err, domain.ErrPollNotFound) {
			return ps.bot.SendMessage("Poll not found")
		}
		log.Printf("GetResults error: %v", err)
		return ps.bot.SendMessage("Failed to get poll results")
	}

	poll, err := ps.repo.GetPoll(pollID)
	if err != nil {
		log.Printf("GetPoll error: %v", err)
		return ps.bot.SendMessage("Failed to get poll info")
	}

	message := fmt.Sprintf("Results for poll: %s\n", poll.Question)
	for optID, count := range results {
		if optText, ok := poll.Options[optID]; ok {
			message += fmt.Sprintf("- %s: %d votes\n", optText, count)
		}
	}

	return ps.bot.SendMessage(message)
}

func (ps *PollService) handleEndPoll(userID, args string) error {
	pollID := strings.TrimSpace(args)
	if pollID == "" {
		return ps.bot.SendMessage("Please specify poll ID")
	}

	if err := ps.repo.EndPoll(pollID, userID); err != nil {
		if errors.Is(err, domain.ErrPollNotFound) {
			return ps.bot.SendMessage("Poll not found")
		}
		if errors.Is(err, domain.ErrNotPollCreator) {
			return ps.bot.SendMessage("Only the poll creator can end the poll")
		}
		log.Printf("EndPoll error: %v", err)
		return ps.bot.SendMessage("Failed to end poll")
	}

	return ps.bot.SendMessage("Poll has been ended successfully")
}

func (ps *PollService) handleDeletePoll(userID, args string) error {
	pollID := strings.TrimSpace(args)
	if pollID == "" {
		return ps.bot.SendMessage("Please specify poll ID")
	}

	if err := ps.repo.DeletePoll(pollID, userID); err != nil {
		if errors.Is(err, domain.ErrPollNotFound) {
			return ps.bot.SendMessage("Poll not found")
		}
		if errors.Is(err, domain.ErrNotPollCreator) {
			return ps.bot.SendMessage("Only the poll creator can delete the poll")
		}
		log.Printf("DeletePoll error: %v", err)
		return ps.bot.SendMessage("Failed to delete poll")
	}

	return ps.bot.SendMessage("Poll has been deleted successfully")
}
