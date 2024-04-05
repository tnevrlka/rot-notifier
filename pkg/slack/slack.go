package slack

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
)

type User struct {
	GitHubUsername string `json:"username"`
	SlackId        string `json:"id"`
}

type Service struct {
	Client *slack.Client
	Users  []User
}

func usersFromBase64JSON(encoded string) ([]User, error) {
	jsonFile, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64: %+v", err)
	}
	var users []User
	err = json.Unmarshal(jsonFile, &users)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %+v", err)
	}
	return users, nil
}

func NewService(token, usernameBase64 string) (*Service, error) {
	client := slack.New(token)
	users, err := usersFromBase64JSON(usernameBase64)
	if err != nil {
		return nil, err
	}
	return &Service{
		Client: client,
		Users:  users,
	}, nil
}

func (service *Service) SlackIdFromGitHubUsername(username string) string {
	for _, user := range service.Users {
		if user.GitHubUsername == username {
			return user.SlackId
		}
	}
	return ""
}

func (service *Service) SendMessage(username string, options ...slack.MsgOption) error {
	channelId := service.SlackIdFromGitHubUsername(username)
	if channelId == "" {
		return fmt.Errorf("user with username '%s' was not found", username)
	}
	_, _, err := service.Client.PostMessage(channelId, options...)
	if err != nil {
		return fmt.Errorf("could not post message: %+v", err)
	}
	return nil
}
