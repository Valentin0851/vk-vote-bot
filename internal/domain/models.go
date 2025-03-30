package domain

type Poll struct {
	ID        string
	CreatorID string
	Question  string
	Options   map[string]string
	CreatedAt int64
	IsActive  bool
}

type PollRepository interface {
	CreatePoll(poll Poll) (string, error)
	Vote(pollID, userID, option string) error
	GetResults(pollID string) (map[string]int, error)
	EndPoll(pollID, userID string) error
	DeletePoll(pollID, userID string) error
}
