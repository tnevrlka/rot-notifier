package types

import "github.com/tnevrlka/rot-notifier/pkg/github"

type MessageInterface interface {
	SendMessage(username, message string) error
}
type Notifier struct {
	GitHubService  *github.Service
	MessageService MessageInterface
}
