package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v60/github"
)

type IssueInterface interface {
	ListComments(ctx context.Context, owner, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error)
	ListIssueEvents(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.IssueEvent, *github.Response, error)
}

type PullRequestInterface interface {
	ListReviews(ctx context.Context, owner, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
	List(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

type Service struct {
	Issues       IssueInterface
	PullRequests PullRequestInterface
	Owner        string
	Repository   string
}

type RepositoryError struct {
	*Service
	Message      string
	WrappedError error
}

type PullRequestError struct {
	*Service
	Message      string
	Number       int
	WrappedError error
}

func (r *RepositoryError) Error() string {
	return fmt.Sprintf("%s in %s/%s: %+v", r.Message, r.Owner, r.Repository, r.WrappedError)
}

func (p *PullRequestError) Error() string {
	return fmt.Sprintf("%s in %s/%s, number %d: %+v", p.Message, p.Owner, p.Repository, p.Number, p.WrappedError)
}

type Assignee struct {
	*github.User
	RequestedAt  github.Timestamp
	LatestAction github.Timestamp
}

type PullRequest struct {
	Number    int
	Owner     *Assignee
	Reviewers []*Assignee
}
