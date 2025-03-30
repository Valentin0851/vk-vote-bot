package bot

import (
	"encoding/json"
	"log"
	"sync"

	"mattermost-vote-bot/internal/domain"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

type Bot struct {
	client      *model.Client4
	wsClient    *model.WebSocketClient
	userID      string
	channelID   string
	pollService *PollService
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func NewBot(cfg *domain.MattermostConfig, repo domain.PollRepository) (*Bot, error) {
	client := model.NewAPIv4Client(cfg.URL)
	client.SetToken(cfg.Token)

	user, _, err := client.GetUserByUsername(cfg.Username, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bot user")
	}

	wsClient, err := model.NewWebSocketClient4(cfg.URL, cfg.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create websocket client")
	}

	return &Bot{
		client:      client,
		wsClient:    wsClient,
		userID:      user.Id,
		channelID:   cfg.Channel,
		pollService: NewPollService(repo),
		stopChan:    make(chan struct{}),
	}, nil
}

func (b *Bot) Start() error {
	b.wsClient.Listen()

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case event := <-b.wsClient.EventChannel:
				if err := b.handleWebSocketEvent(event); err != nil {
					log.Printf("Error handling event: %v", err)
				}
			case <-b.stopChan:
				return
			}
		}
	}()

	log.Println("Bot started successfully")
	return nil
}

func (b *Bot) Stop() {
	close(b.stopChan)
	b.wg.Wait()
	b.wsClient.Close()
	log.Println("Bot stopped gracefully")
}

func (b *Bot) SendMessage(message string) error {
	post := &model.Post{
		ChannelId: b.channelID,
		Message:   message,
	}

	_, _, err := b.client.CreatePost(post)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}
	return nil
}

func (b *Bot) handleWebSocketEvent(event *model.WebSocketEvent) error {
	if event.EventType() != model.WebsocketEventPosted {
		return nil
	}

	var post model.Post
	if err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post); err != nil {
		return errors.Wrap(err, "failed to unmarshal post")
	}

	if post.ChannelId != b.channelID || post.UserId == b.userID {
		return nil
	}

	return b.pollService.HandleCommand(post)
}
