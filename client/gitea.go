package client

import (
	"errors"
	"strings"

	sdk "code.gitea.io/sdk/gitea"
)

type gitea struct {
	endpoint  string
	token     string
	isOrg     bool
	isPrivate bool

	client *sdk.Client
}

var _ IClient = (*gitea)(nil)

func NewGiteaClient(endpoint string, token string, isOrg, isPrivate bool) (*gitea, error) {
	var opts = []sdk.ClientOption{
		sdk.SetToken(token),
	}
	client, err := sdk.NewClient(endpoint, opts...)
	if err != nil {
		return nil, err
	}
	return &gitea{endpoint: endpoint, token: token, isOrg: isOrg, isPrivate: isPrivate, client: client}, nil
}

func (g *gitea) IsRepositoryExist(owner, repository string) (bool, error) {
	repo, _, err := g.client.GetRepo(owner, repository)
	if err != nil {
		if strings.Contains(err.Error(), "The target couldn't be found.") {
			return false, nil
		}
		return false, err
	}
	if repo == nil {
		return false, nil
	}
	return true, nil
}

func (g *gitea) CreateRepository(owner, repository string) error {
	if g.isOrg {
		repo, _, err := g.client.CreateOrgRepo(owner, sdk.CreateRepoOption{
			Name:          repository,
			Private:       g.isPrivate,
			DefaultBranch: "main",
		})
		if err != nil {
			return err
		}
		if repo == nil {
			return errors.New("create organization repository failed: response empty")
		}
		return nil
	}
	repo, _, err := g.client.CreateRepo(sdk.CreateRepoOption{
		Name:          repository,
		Private:       g.isPrivate,
		DefaultBranch: "main",
	})
	if err != nil {
		return err
	}
	if repo == nil {
		return errors.New("create repository failed: response empty")
	}
	return nil
}
