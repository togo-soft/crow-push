package client

import (
	"strings"

	sdk "github.com/ktrysmt/go-bitbucket"
)

type bitbucket struct {
	username, password string
	isOrg, isPrivate   bool

	client *sdk.Client
}

func NewBitBucketClient(username, password string, isOrg, isPrivate bool) (*bitbucket, error) {
	client, err := sdk.NewBasicAuth(username, password)
	if err != nil {
		return nil, err
	}
	return &bitbucket{
		username:  username,
		password:  password,
		isOrg:     isOrg,
		isPrivate: isPrivate,
		client:    client,
	}, nil
}

var _ IClient = (*bitbucket)(nil)

func (b *bitbucket) IsRepositoryExist(owner, repository string) (bool, error) {
	repo, err := b.client.Repositories.Repository.Get(&sdk.RepositoryOptions{
		Owner:    owner,
		RepoSlug: repository,
		Scm:      "git",
	})
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		}
		return false, err
	}
	if repo == nil {
		return false, nil
	}
	return true, nil
}

func (b *bitbucket) CreateRepository(owner, repository string) error {
	isPrivate := "false"
	if b.isPrivate {
		isPrivate = "true"
	}
	_, err := b.client.Repositories.Repository.Create(&sdk.RepositoryOptions{
		Owner:     owner,
		RepoSlug:  repository,
		Scm:       "git",
		IsPrivate: isPrivate,
	})
	if err != nil {
		return err
	}
	return nil
}
