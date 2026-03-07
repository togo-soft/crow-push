package client

import (
	"testing"
)

func TestNewGithubClient(t *testing.T) {
	client := NewGithubClient("token", false, false)

	exist, err := client.IsRepositoryExist("owner name", "repo name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
