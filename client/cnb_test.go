package client

import (
	"testing"
)

func TestNewCNBClient(t *testing.T) {
	client, err := NewCNBClient("token", false, false)
	if err != nil {
		t.Fatal(err)
	}
	exist, err := client.IsRepositoryExist("owner name", "repo name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}
