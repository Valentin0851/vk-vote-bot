package mattermost

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type Client struct {
	API       *model.Client4
	WebSocket *model.WebSocketClient
	UserID    string
	ChannelID string
}

func NewClient(url, token, team, channel string) (*Client, error) {
	api := model.NewAPIv4Client(url)
	api.SetToken(token)

	user, _, err := api.GetUserByUsername(team, "")
	if err != nil {
		return nil, err
	}

	wsClient, err := model.NewWebSocketClient4(url, token)
	if err != nil {
		return nil, err
	}

	return &Client{
		API:       api,
		WebSocket: wsClient,
		UserID:    user.Id,
		ChannelID: channel,
	}, nil
}
