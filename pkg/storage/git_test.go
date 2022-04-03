package storage

import (
	"context"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"testing"
)

const repo = "https://github.com/neoboxer/outbox-config.git"

func TestNewGit(t *testing.T) {
	auth := &http.BasicAuth{
		Username: "ryan-ovo",
		Password: "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",
	}
	g := NewGit(repo, NewGitOption().WithAuth(auth))
	ctx := context.Background()
	if _, err := g.Env(ctx, "test"); err != nil {
		t.Fatal(err)
	}
	file, err := g.Use(ctx, "default")
	if err != nil {
		t.Fatal(err)
	}
	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
