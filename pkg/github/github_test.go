package github

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v60/github"
	"github.com/stretchr/testify/mock"
)

type MockPullRequest struct {
	mock.Mock
}

type MockIssues struct {
	mock.Mock
}

func (m *MockIssues) ListComments(ctx context.Context, owner, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
	args := m.Called(ctx, owner, repo, number, opts)
	return args.Get(0).([]*github.IssueComment), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockIssues) ListIssueEvents(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.IssueEvent, *github.Response, error) {
	args := m.Called(ctx, owner, repo, number, opts)
	return args.Get(0).([]*github.IssueEvent), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockPullRequest) ListReviews(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
	args := m.Called(ctx, owner, repo, number, opts)
	return args.Get(0).([]*github.PullRequestReview), args.Get(1).(*github.Response), args.Error(2)
}

func (m *MockPullRequest) List(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	args := m.Called(ctx, owner, repo, opts)
	return args.Get(0).([]*github.PullRequest), args.Get(1).(*github.Response), args.Error(2)
}

// TestNewService
func TestNewService(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		expected    *Service
		valid       bool
	}{
		{
			name:        "authoriezed",
			accessToken: "test-token",
			expected: &Service{
				Owner:      "redhat-appstudio",
				Repository: "e2e-tests",
			},
			valid: true,
		},
		{
			name:        "unauthorized",
			accessToken: "",
			expected: &Service{
				Owner:      "redhat-appstudio",
				Repository: "e2e-tests",
			},
			valid: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewService("redhat-appstudio", "e2e-tests", tc.accessToken)
			if !tc.valid && service != nil {
				t.Errorf("expected nil, got %v", service)
			}
			if tc.valid {
				if service.Owner != tc.expected.Owner {
					t.Errorf("expected owner '%s', got '%s'", tc.expected.Owner, service.Owner)
				}
				if service.Repository != tc.expected.Repository {
					t.Errorf("expected repository '%s', got '%s'", tc.expected.Repository, service.Repository)
				}
			}
			if service.PullRequests == nil {
				t.Error("expected PullRequests to be set")
			}
			if service.Issues == nil {
				t.Error("expected Issues to be set")
			}
		})
	}
}

func TestListIssueComments(t *testing.T) {
	testCases := []struct {
		name    string
		ghError error
	}{
		{
			name:    "valid",
			ghError: nil,
		},
		{
			name:    "error",
			ghError: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockIssues := new(MockIssues)
			service := Service{
				Issues:       mockIssues,
				PullRequests: &MockPullRequest{},
				Owner:        "redhat-appstudio",
				Repository:   "e2e-tests",
			}
			mockIssues.On("ListComments", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*github.IssueComment{}, &github.Response{}, tc.ghError)
			_, err := service.ListIssueComments(1)
			var pullRequestError error = nil
			if tc.ghError != nil {
				pullRequestError = &PullRequestError{
					Service:      &service,
					Message:      "error listing issue comments",
					Number:       1,
					WrappedError: tc.ghError,
				}
			}
			if err != nil && pullRequestError != nil && err.Error() != pullRequestError.Error() {
				t.Errorf("expected error '%v', got '%v'", pullRequestError.Error(), err)
			}
			mockIssues.AssertExpectations(t)
		})
	}
}

func TestListPullRequestReviews(t *testing.T) {
	testCases := []struct {
		name    string
		ghError error
	}{
		{
			name:    "valid",
			ghError: nil,
		},
		{
			name:    "error",
			ghError: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPullRequest := new(MockPullRequest)
			service := Service{
				Issues:       &MockIssues{},
				PullRequests: mockPullRequest,
				Owner:        "redhat-appstudio",
				Repository:   "e2e-tests",
			}
			mockPullRequest.On("ListReviews", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*github.PullRequestReview{}, &github.Response{}, tc.ghError)
			_, err := service.ListPullRequestReviews(1)
			var pullRequestError error = nil
			if tc.ghError != nil {
				pullRequestError = &PullRequestError{
					Service:      &service,
					Message:      "error listing pull request reviews",
					Number:       1,
					WrappedError: tc.ghError,
				}
			}
			if err != nil && pullRequestError != nil && err.Error() != pullRequestError.Error() {
				t.Errorf("expected error '%v', got '%v'", pullRequestError.Error(), err)
			}
			mockPullRequest.AssertExpectations(t)
		})
	}
}

func TestListPullRequestEvents(t *testing.T) {
	testCases := []struct {
		name    string
		ghError error
	}{
		{
			name:    "valid",
			ghError: nil,
		},
		{
			name:    "error",
			ghError: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockIssues := new(MockIssues)
			service := Service{
				Issues:       mockIssues,
				PullRequests: &MockPullRequest{},
				Owner:        "redhat-appstudio",
				Repository:   "e2e-tests",
			}
			mockIssues.On("ListIssueEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*github.IssueEvent{}, &github.Response{}, tc.ghError)
			_, err := service.ListPullRequestEvents(1)
			var pullRequestError error = nil
			if tc.ghError != nil {
				pullRequestError = &PullRequestError{
					Service:      &service,
					Message:      "error listing pull request events",
					Number:       1,
					WrappedError: tc.ghError,
				}
			}
			if err != nil && pullRequestError != nil && err.Error() != pullRequestError.Error() {
				t.Errorf("expected error '%v', got '%v'", pullRequestError.Error(), err)
			}
			mockIssues.AssertExpectations(t)
		})
	}
}

func TestListPullRequestReviewRequests(t *testing.T) {
	testCases := []struct {
		name    string
		events  []*github.IssueEvent
		ghError error
	}{
		{
			name: "valid",
			// write a bunch of events
			events: []*github.IssueEvent{
				{
					Event: github.String("review_requested"),
				},
				{
					Event: github.String("foo"),
				},
			},
			ghError: nil,
		},
		{
			name:    "error",
			ghError: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockIssues := new(MockIssues)
			service := Service{
				Issues:       mockIssues,
				PullRequests: &MockPullRequest{},
				Owner:        "redhat-appstudio",
				Repository:   "e2e-tests",
			}
			mockIssues.On("ListIssueEvents", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.events, &github.Response{}, tc.ghError)
			reviewRequests, err := service.ListPullRequestReviewRequests(1)
			var pullRequestError error = nil
			if tc.ghError != nil {
				pullRequestError = &PullRequestError{
					Service:      &service,
					Message:      "error listing pull request events",
					Number:       1,
					WrappedError: tc.ghError,
				}
			}
			if err != nil && pullRequestError != nil && err.Error() != pullRequestError.Error() {
				t.Errorf("expected error '%v', got '%v'", pullRequestError.Error(), err)
			}
			for _, event := range reviewRequests {
				if event.GetEvent() != "review_requested" {
					t.Errorf("expected 'review_requested', got '%s'", event.GetEvent())
				}
			}
			mockIssues.AssertExpectations(t)
		})
	}
}

func TestListOpenPullRequests(t *testing.T) {
	testCases := []struct {
		name    string
		ghError error
	}{
		{
			name:    "valid",
			ghError: nil,
		},
		{
			name:    "error",
			ghError: errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPullRequest := new(MockPullRequest)
			service := Service{
				Issues:       &MockIssues{},
				PullRequests: mockPullRequest,
				Owner:        "redhat-appstudio",
				Repository:   "e2e-tests",
			}
			mockPullRequest.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*github.PullRequest{}, &github.Response{}, tc.ghError)
			_, err := service.ListOpenPullRequests()
			var pullRequestError error = nil
			if tc.ghError != nil {
				pullRequestError = &RepositoryError{
					Service:      &service,
					Message:      "error listing open pull requests",
					WrappedError: tc.ghError,
				}
			}
			if err != nil && pullRequestError != nil && err.Error() != pullRequestError.Error() {
				t.Errorf("expected error '%v', got '%v'", pullRequestError.Error(), err)
			}
			mockPullRequest.AssertExpectations(t)
		})
	}
}
