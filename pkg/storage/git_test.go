package storage

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

const repo = "https://github.com/neoboxer/outbox-config.git"
const username = "access_token_with_any_non_empty_string"
const accessToken = "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"

func TestNewGit(t *testing.T) {
	ctx := context.Background()

	auth := &http.BasicAuth{
		Username: username,
		Password: accessToken,
	}
	g := NewGit(repo, NewGitOption().WithAuth(auth).WithEnv("test"))

	file, err := g.Use(ctx, "default")
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
