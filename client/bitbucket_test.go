package client

import (
	"testing"

	sdk "github.com/ktrysmt/go-bitbucket"
)

func TestNewBitBucketClient(t *testing.T) {
	client, err := sdk.NewBasicAuth("your_email", "bitbucket_scope_api_token")
	if err != nil {
		t.Fatal(err)
	}
	get, err := client.Repositories.Repository.Get(&sdk.RepositoryOptions{
		Owner:    "workspace name",
		RepoSlug: "repository name",
		Scm:      "git",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(get)
}
