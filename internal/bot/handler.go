package bot

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

func (ps *PollService) HandleCommand(post model.Post) error {
	if !strings.HasPrefix(post.Message, "/poll") {
		return nil
	}

	args := strings.Fields(post.Message)
	if len(args) < 2 {
		return ps.bot.SendMessage("Invalid command format. Usage: /poll [create|vote|results|end|delete] [args]")
	}

	command := args[1]
	rest := strings.Join(args[2:], " ")

	switch command {
	case "create":
		return ps.handleCreatePoll(post.UserId, post.ChannelId, rest)
	case "vote":
		return ps.handleVote(post.UserId, rest)
	case "results":
		return ps.handleResults(rest)
	case "end":
		return ps.handleEndPoll(post.UserId, rest)
	case "delete":
		return ps.handleDeletePoll(post.UserId, rest)
	default:
		return ps.bot.SendMessage(fmt.Sprintf("Unknown command: %s. Available commands: create, vote, results, end, delete", command))
	}
}
