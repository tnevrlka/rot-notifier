package slack

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestUsersFromBase64JSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []User
		valid    bool
	}{
		{
			name: "valid json",
			input: []byte(
				`[
	{"username": "foo", "id": "U123"},
	{"username": "spam", "id": "UUUU"},
	{"username": "user", "id": "UKYR"}
]`),

			expected: []User{
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
			},
			valid: true,
		},
		{
			name: "invalid json",
			input: []byte(
				`[
	{"username": "foo", "id": "U123"},
	{"username": "spam", "id": "UUUU"},
	{"username": "user", "id": "UKYR"},
]`),
			expected: nil,
			valid:    false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			users, err := usersFromBase64JSON(base64.StdEncoding.EncodeToString(testCase.input))
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
	validService, err := NewService("foo",
		"WyAKCXsidXNlcm5hbWUiOiAiZm9vIiwgImlkIjogIlUxMjMifSwKCXsidXNlcm5hbWUiOiAic3BhbSIsICJpZCI6ICJVVVVVIn0sCgl7InVzZXJuYW1lIjogInVzZXIiLCAiaWQiOiAiVUtZUiJ9Cl0K")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name     string
		service  *Service
		input    string
		expected string
	}{
		{
			name:     "user exists",
			service:  validService,
			input:    "foo",
			expected: "U123",
		},
		{
			name:     "user does not exist",
			service:  validService,
			input:    "ferda",
			expected: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			slackId := testCase.service.SlackIdFromGitHubUsername(testCase.input)
			if slackId != testCase.expected {
				t.Errorf("expected slackId of '%s' to be '%s', got '%s'", testCase.input, testCase.expected, slackId)
			}
		})
	}
}
