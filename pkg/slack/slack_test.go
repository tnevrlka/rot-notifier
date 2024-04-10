package slack

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

var sampleUsers = []User{
	{
		GitHubUsername: "foo",
		SlackId:        "U123",
	},
	{
		GitHubUsername: "spam",
		SlackId:        "UUUU",
	},
	{
		GitHubUsername: "user",
		SlackId:        "UKYR",
	},
}

var jsonBytes, _ = json.Marshal(sampleUsers)
var base64File = base64.StdEncoding.EncodeToString(jsonBytes)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	args := m.Called(channelID, options)
	return args.String(0), args.String(1), args.Error(2)
}

func TestNewService(t *testing.T) {
	testCases := []struct {
		name           string
		token          string
		usernameBase64 string
		expected       *Service
		valid          bool
	}{
		{
			name:           "valid",
			usernameBase64: base64File,
			expected: &Service{
				Users: sampleUsers,
			},
			valid: true,
		},
		{
			name:           "invalid",
			usernameBase64: base64File + "foo",
			expected:       nil,
			valid:          false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service, err := NewService("", testCase.usernameBase64)
			if !testCase.valid && err == nil {
				t.Errorf("expected err to occurr, was: %v", err)
			}
			if testCase.valid {
				if err != nil {
					t.Errorf("expected err to not occurr: %v", err)
				}
				// Check that users are the same
				for _, expUser := range testCase.expected.Users {
					found := false
					for _, user := range service.Users {
						if user.GitHubUsername == expUser.GitHubUsername && user.SlackId == expUser.SlackId {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("could not find expected user '%v' in actual users '%v'", expUser, service.Users)
					}
				}
			}
		})
	}
}

func TestUsersFromBase64JSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []User
		valid    bool
	}{
		{
			name:     "valid json",
			input:    base64File,
			expected: sampleUsers,
			valid:    true,
		},
		{
			name: "invalid json",
			input: base64.StdEncoding.EncodeToString([]byte(
				`[
	{"username": "foo", "id": "U123"},
	{"username": "spam", "id": "UUUU"},
	{"username": "user", "id": "UKYR"},
]`)),
			expected: nil,
			valid:    false,
		},
		{
			name:     "invalid base64",
			input:    base64File + "foo",
			expected: nil,
			valid:    false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			users, err := usersFromBase64JSON(testCase.input)
			if testCase.valid && err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(users, testCase.expected) {
				t.Errorf("expected '%v', got '%v'", testCase.expected, users)
			}
		})
	}
}

func TestSlackIdFromGitHubUsername(t *testing.T) {
	mockService := Service{Users: sampleUsers}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "user exists",
			input:    sampleUsers[0].GitHubUsername,
			expected: sampleUsers[0].SlackId,
		},
		{
			name:     "user does not exist",
			input:    "ferda",
			expected: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			slackId := mockService.MsgIdFromGitUsername(testCase.input)
			if slackId != testCase.expected {
				t.Errorf("expected slackId of '%s' to be '%s', got '%s'", testCase.input, testCase.expected, slackId)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	mockService := Service{Users: sampleUsers}

	testCases := []struct {
		name        string
		username    string
		channelId   string
		options     []slack.MsgOption
		expected    error
		postError   error
		reachesPost bool
	}{
		{
			name:      "valid",
			username:  mockService.Users[0].GitHubUsername,
			channelId: mockService.Users[0].SlackId,
			options: []slack.MsgOption{
				slack.MsgOptionText("foo", false),
				slack.MsgOptionAsUser(true),
			},
			expected:    nil,
			postError:   nil,
			reachesPost: true,
		},
		{
			name:      "invalid username",
			username:  "ferda",
			channelId: "",
			options: []slack.MsgOption{
				slack.MsgOptionText("foo", false),
				slack.MsgOptionAsUser(true),
			},
			expected:    errors.New("user with username 'ferda' was not found"),
			postError:   nil,
			reachesPost: false,
		},
		{
			name:      "post error",
			username:  mockService.Users[0].GitHubUsername,
			channelId: mockService.Users[0].SlackId,
			options: []slack.MsgOption{
				slack.MsgOptionText("foo", false),
				slack.MsgOptionAsUser(true),
			},
			expected:    errors.New("could not post message: test post error"),
			postError:   errors.New("test post error"),
			reachesPost: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockClient := MockClient{}
			service := Service{Users: mockService.Users, Client: &mockClient}

			if testCase.reachesPost {
				mockClient.On("PostMessage", testCase.channelId, mock.Anything).Return("", "", testCase.postError)
			}
			err := service.SendMessage(testCase.username, "foo")
			if err != nil && testCase.expected != nil && err.Error() != testCase.expected.Error() {
				t.Errorf("errors do not match, got '%v', expected '%v'", err, testCase.expected)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
