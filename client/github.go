package client

import (
	"context"
	"errors"
	"strings"

	sdk "github.com/google/go-github/v84/github"
)

type github struct {
	token     string
	isOrg     bool
	isPrivate bool

	ctx    context.Context
	client *sdk.Client
}

var _ IClient = (*github)(nil)

func NewGithubClient(token string, isOrg, isPrivate bool) *github {
	client := sdk.NewClient(nil).WithAuthToken(token)
	return &github{token: token, isOrg: isOrg, isPrivate: isPrivate, ctx: context.Background(), client: client}
}

func (g *github) IsRepositoryExist(owner, repository string) (bool, error) {
	get, _, err := g.client.Repositories.Get(g.ctx, owner, repository)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		}
		return false, err
	}
	if get == nil {
		return false, nil
	}
	return true, nil
}

func (g *github) CreateRepository(owner, repository string) error {
	org := ""
	if g.isOrg {
		org = owner
	}
	create, _, err := g.client.Repositories.Create(g.ctx, org, &sdk.Repository{
		Name:              &repository,
		DefaultBranch:     new("main"),
		MasterBranch:      new("main"),
		AutoInit:          new(false),
		AllowRebaseMerge:  new(true),
		AllowUpdateBranch: new(true),
		AllowSquashMerge:  new(true),
		AllowMergeCommit:  new(true),
		AllowAutoMerge:    new(true),
		AllowForking:      new(true),
		Private:           &g.isPrivate,
		HasIssues:         new(true),
	})
	if err != nil {
		return err
	}
	if create == nil {
		return errors.New("create repository failed: response is nil")
	}
	return nil
}
