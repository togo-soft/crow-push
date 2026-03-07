package client

import (
	"testing"

	"codefloe.com/actions/common"
)

func TestNewGiteaClient(t *testing.T) {
	codefloe, err := NewGiteaClient(common.DefaultEndpoint, "gitea platform token", true, false)
	if err != nil {
		t.Fatal(err)
	}
	exist, err := codefloe.IsRepositoryExist("owner name", "repo name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
