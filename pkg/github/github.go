package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	"net/http"
)

type Service struct {
	Client     *github.Client
	Owner      string
	Repository string
}

type RepositoryError struct {
	Service
	Message string
}

type PullRequestError struct {
	Service
	Message string
	Number  int
}

func (r *RepositoryError) Error() string {
	return fmt.Sprintf("%s in %s/%s", r.Message, r.Owner, r.Repository)
}

func (p *PullRequestError) Error() string {
	return fmt.Sprintf("%s in %s/%s, number %d", p.Message, p.Owner, p.Repository, p.Number)
}

func NewService(owner, repository, accessToken string) Service {
	var tc *http.Client = nil
	if accessToken != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	client := github.NewClient(tc)
	return Service{
		Client:     client,
		Owner:      owner,
		Repository: repository,
	}
}

func (service *Service) ListIssueComments(prNumber int) ([]*github.IssueComment, error) {
	comments, _, err := service.Client.Issues.ListComments(context.TODO(), service.Owner, service.Repository, prNumber, nil)
	if err != nil {
		return nil, &PullRequestError{
			Service: *service,
			Message: "error listing issue comments",
			Number:  prNumber,
		}
	}
	return comments, nil
}

func (service *Service) ListPullRequestReviews(prNumber int) ([]*github.PullRequestReview, error) {
	reviews, _, err := service.Client.PullRequests.ListReviews(context.TODO(), service.Owner, service.Repository, prNumber, nil)
	if err != nil {
		return nil, &PullRequestError{
			Service: *service,
			Message: "error listing pull request reviews",
			Number:  prNumber,
		}
	}
	return reviews, nil
}

func (service *Service) ListPullRequestEvents(prNumber int) ([]*github.IssueEvent, error) {
	events, _, err := service.Client.Issues.ListIssueEvents(context.TODO(), service.Owner, service.Repository, prNumber, nil)
	if err != nil {
		return nil, &PullRequestError{
			Service: *service,
			Message: "error listing pull request events",
			Number:  prNumber,
		}
	}
	return events, nil
}

func (service *Service) ListPullRequestReviewRequests(prNumber int) ([]*github.IssueEvent, error) {
	events, err := service.ListPullRequestEvents(prNumber)
	if err != nil {
		return nil, err
	}

	var requestedReviews []*github.IssueEvent
	for _, event := range events {
		if event.GetEvent() == "review_requested" {
			requestedReviews = append(requestedReviews, event)
		}
	}
	return requestedReviews, nil
}

func (service *Service) ListOpenPullRequests() ([]*github.PullRequest, error) {
	list, _, err := service.Client.PullRequests.List(context.TODO(), service.Owner, service.Repository, nil)
	if err != nil {
		return nil, &RepositoryError{
			Service: *service,
			Message: "error listing open pull requests",
		}
	}
	return list, nil
}
