package client

import (
	"testing"
)

func TestNewGiteeClient(t *testing.T) {
	cli := NewGiteeClient("token", false, false)

	exist, err := cli.IsRepositoryExist("owner name", "repo name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
