package client

import (
	"context"

	sdk "cnb.cool/cnb/sdk/go-cnb/cnb"
	"cnb.cool/cnb/sdk/go-cnb/cnb/types/dto"
)

type cnb struct {
	token     string
	isOrg     bool
	isPrivate bool

	ctx    context.Context
	client *sdk.Client
}

func NewCNBClient(token string, isOrg, isPrivate bool) (*cnb, error) {
	client, err := sdk.NewClient(nil).WithAuthToken(token).WithURLs("https://api.cnb.cool/")
	if err != nil {
		return nil, err
	}
	return &cnb{ctx: context.Background(), client: client, token: token, isOrg: isOrg, isPrivate: isPrivate}, nil
}

var _ IClient = (*cnb)(nil)

func (c *cnb) IsRepositoryExist(owner, repository string) (bool, error) {
	repos, _, err := c.client.Repositories.GetRepos(c.ctx, &sdk.GetReposOptions{
		Search: repository,
	})
	if err != nil {
		return false, err
	}
	excepted := owner + "/" + repository
	for _, repo := range repos {
		if repo.Name != repository {
			continue
		}
		if excepted != repo.Path {
			continue
		}
		return true, nil
	}
	return false, nil
}

func (c *cnb) CreateRepository(owner, repository string) error {
	visibility := "public"
	if c.isPrivate {
		visibility = "private"
	}
	_, err := c.client.Repositories.CreateRepo(c.ctx, owner, &sdk.CreateRepoRequest{
		Name:       repository,
		Visibility: dto.CreateRepoReqVisibility(visibility),
	})
	if err != nil {
		return err
	}
	return nil
}
