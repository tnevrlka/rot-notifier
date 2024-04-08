package slack

import (
	"encoding/base64"
	"reflect"
	"testing"
)

// Base64 encoded sample valid JSON
// `[
//
//	{"username": "foo", "id": "U123"},
//	{"username": "spam", "id": "UUUU"},
//	{"username": "user", "id": "UKYR"}
//
// ]`)
const base64File = "WyAKCXsidXNlcm5hbWUiOiAiZm9vIiwgImlkIjogIlUxMjMifSwKCXsidXNlcm5hbWUiOiAic3BhbSIsICJpZCI6ICJVVVVVIn0sCgl7InVzZXJuYW1lIjogInVzZXIiLCAiaWQiOiAiVUtZUiJ9Cl0K"

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

func NewMockService(base64String string) (*Service, error) {
	service := new(Service)
	var err error
	service.Users, err = usersFromBase64JSON(base64String)
	return service, err
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
	mockService, err := NewMockService(base64File)
	if err != nil {
		t.Fatal(err)
	}
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
			slackId := mockService.SlackIdFromGitHubUsername(testCase.input)
			if slackId != testCase.expected {
				t.Errorf("expected slackId of '%s' to be '%s', got '%s'", testCase.input, testCase.expected, slackId)
			}
		})
	}
}
