package client

import (
	"errors"

	sdk "gitee.com/sdk/golang-sdk-v5/client"
	"gitee.com/sdk/golang-sdk-v5/client/repositories"
)

type gitee struct {
	token     string
	isOrg     bool
	isPrivate bool
}

func NewGiteeClient(token string, isOrg, isPrivate bool) *gitee {
	return &gitee{token: token, isOrg: isOrg, isPrivate: isPrivate}
}

func (g *gitee) IsRepositoryExist(owner, repository string) (bool, error) {
	repos, err := sdk.Default.Repositories.GetV5UserRepos(&repositories.GetV5UserReposParams{
		AccessToken: &g.token,
		Q:           &repository,
	})
	// request failed
	if err != nil {
		return false, err
	}
	// no repo list
	if len(repos.Payload) == 0 {
		return false, nil
	}
	for _, project := range repos.Payload {
		if project.Name != repository {
			continue
		}
		if project.Namespace.Path != owner {
			continue
		}
		return true, nil
	}
	return false, nil
}

func (g *gitee) CreateRepository(owner, repository string) error {
	// default public is true, 0 private, 1 public, 2 org public
	var public int32 = 1
	if g.isPrivate {
		public = 0
	}
	if g.isOrg {
		response, err := sdk.Default.Repositories.CreateOrgOrgRepos(&repositories.CreateOrgOrgReposParams{
			AccessToken: &g.token,
			Name:        repository,
			Org:         owner,
			Public:      &public,
		})
		if err != nil {
			return err
		}
		if response.GetPayload() == nil {
			return errors.New("create organization repository failed: payload is nil")
		}
		return nil
	}
	response, err := sdk.Default.Repositories.PostV5UserRepos(&repositories.PostV5UserReposParams{})
	if err != nil {
		return err
	}
	if response.GetPayload() == nil {
		return errors.New("create user repository failed: payload is nil")
	}
	return nil
}

var _ IClient = (*gitee)(nil)
