package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	defaultFetchTimeout = 2 * time.Second
	defaultFetchTTL     = 10 * time.Second
	defaultFilePrefix   = "jsonnet"
)

type gitInfo struct {
	fs     internalFs
	remote plumbing.Hash
	local  plumbing.Hash
}

func (g gitInfo) isSameHash() bool {
	// compare local and remote is consistent
	return g.local == g.remote
}

type Git struct {
	// URL
	URL string
	// FetchTimeout
	FetchTimeout time.Duration
	// FetchTTL
	FetchTTL time.Duration
	// auth
	auth transport.AuthMethod
	// init represents the first fetch operation, remain false until first successful fetch
	init bool
	// infos represents git reference information
	infos map[plumbing.ReferenceName]gitInfo
	// remote
	remote *Remote
	// storage
	storage  storage.Storer
	syncTime time.Time
	lock     sync.Locker
	env      string
	prefix   string
}

type GitOption struct {
	Env          string
	Name         string
	Auth         http.AuthMethod
	FilePrefix   string
	FetchTTL     time.Duration
	FetchTimeout time.Duration
}

func NewGitOption() *GitOption {
	return &GitOption{
		Env:          "master",
		Name:         "origin",
		Auth:         nil,
		FilePrefix:   defaultFilePrefix,
		FetchTTL:     defaultFetchTTL,
		FetchTimeout: defaultFetchTimeout,
	}
}

func (p *GitOption) WithEnv(env string) *GitOption {
	p.Env = env
	return p
}

func (p *GitOption) WithName(name string) *GitOption {
	p.Name = name
	return p
}

func (p *GitOption) WithAuth(auth http.AuthMethod) *GitOption {
	p.Auth = auth
	return p
}

func (p *GitOption) WithFilePrefix(prefix string) *GitOption {
	p.FilePrefix = prefix
	return p
}

func (p *GitOption) WithFetchTTL(TTL time.Duration) *GitOption {
	p.FetchTTL = TTL
	return p
}

func (p *GitOption) WithFetchTimeout(timeout time.Duration) *GitOption {
	p.FetchTimeout = timeout
	return p
}

func NewGit(URL string, options ...*GitOption) *Git {
	var option *GitOption
	if len(options) != 0 {
		option = options[0]
	} else {
		option = NewGitOption()
	}
	store := memory.NewStorage()
	remote := NewRemote(store, &config.RemoteConfig{
		Name:  option.Name,
		URLs:  []string{URL},
		Fetch: []config.RefSpec{"refs/heads/*:refs/heads/*"},
	})
	return &Git{
		URL:          URL,
		auth:         option.Auth,
		FetchTimeout: option.FetchTimeout,
		FetchTTL:     option.FetchTTL,
		infos:        make(map[plumbing.ReferenceName]gitInfo),
		remote:       remote,
		storage:      store,
		lock:         &sync.Mutex{},
		env:          option.Env,
		prefix:       option.FilePrefix,
	}
}

func (g *Git) Env(ctx context.Context, env string) (ReadonlyFs, error) {
	g.lock.Lock()
	g.env = env
	g.lock.Unlock()
	if !g.skipFetch() {
		if err := g.fetch(ctx); err != nil {
			return nil, err
		}
	}

	branch, err := g.getBranch(env)
	if err != nil {
		return nil, err
	}
	if !g.skipCheckout(branch) {
		// switch to target branch fail
		if err = g.checkout(branch); err != nil {
			return nil, err
		}
	}
	return g.infos[branch].fs, nil
}

func (g *Git) Use(ctx context.Context, namespace string) (file ReadonlyFile, err error) {
	var fs ReadonlyFs
	if fs, err = g.Env(ctx, g.env); err != nil {
		return
	}
	file, err = fs.Open(fmt.Sprintf("%s.%s", namespace, g.prefix))
	return
}

func (g *Git) fetch(ctx context.Context) error {
	// lock all the information from Git, keep Git is a correct version.
	g.lock.Lock()
	defer g.lock.Unlock()
	// handle error
	onError := func(err error) error {
		if g.init {
			log.Printf("use old reference data, fetch error = [%+v]\n", err)
			return nil
		}
		log.Printf("fetch error = [%+v]\n", err)
		return err
	}
	timeout := g.FetchTimeout
	if !g.init {
		timeout *= 10
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	err := g.remote.FetchContext(ctx, &FetchOptions{
		RefSpecs: g.remote.Config().Fetch,
		Auth:     g.auth,
		Progress: os.Stderr,
		Force:    true,
	})
	if errors.Is(err, NoErrAlreadyUpToDate) || errors.Is(err, transport.ErrEmptyUploadPackRequest) {
		err = nil
	}
	if err != nil {
		return onError(err)
	}
	ref := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.Master)
	err = g.storage.SetReference(ref)
	if err != nil {
		return onError(err)
	}
	err = g.updateRefs()
	if err != nil {
		return onError(err)
	}
	log.Println("fetch success!")
	g.init = true
	return nil
}

func (g *Git) updateRefs() error {
	it, err := g.storage.IterReferences()
	if err != nil {
		return err
	}
	refs := make(map[plumbing.ReferenceName]plumbing.Hash)
	err = it.ForEach(func(ref *plumbing.Reference) error {
		refs[ref.Name()] = ref.Hash()
		return nil
	})
	if err != nil {
		return err
	}
	next := make(map[plumbing.ReferenceName]gitInfo, len(refs))
	for name, hash := range refs {
		oldInfo := g.infos[name]
		next[name] = gitInfo{
			fs:     oldInfo.fs,
			remote: hash,
			local:  oldInfo.local,
		}
	}
	g.infos = next
	g.syncTime = time.Now()
	return nil
}

func (g *Git) getBranch(env string) (plumbing.ReferenceName, error) {
	refName := g.branchRef(env)
	if _, ok := g.infos[refName]; !ok {
		return "", fmt.Errorf("environment not found, env = [%+v]", env)
	}
	return refName, nil
}

func (g *Git) branchRef(env string) plumbing.ReferenceName {
	return plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", env))
}

func (g *Git) checkout(branch plumbing.ReferenceName) (err error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	var (
		repo *Repository
		tree *Worktree
		head *plumbing.Reference
		fs   internalFs
	)
	fs.Filesystem = memfs.New()
	if repo, err = Open(g.storage, fs.Filesystem); err != nil {
		return
	}
	if tree, err = repo.Worktree(); err != nil {
		return
	}
	if err = tree.Checkout(&CheckoutOptions{
		Branch: branch,
		Force:  true,
	}); err != nil {
		return
	}
	if head, err = repo.Head(); err != nil {
		return
	}
	oldInfo := g.infos[branch]
	g.infos[branch] = gitInfo{
		fs:     fs,
		remote: oldInfo.remote,
		local:  head.Hash(),
	}
	log.Printf("checkout to target branch [%+v] success!\n", branch)
	return
}

func (g *Git) skipFetch() bool {
	return time.Since(g.syncTime) < g.FetchTTL
}

func (g *Git) skipCheckout(branch plumbing.ReferenceName) bool {
	return g.infos[branch].isSameHash()
}

// 对于git文件系统而言，需要提供哪些信息
// 如何表示远程分支 -> 本地分支的对应关系

type internalFs struct {
	billy.Filesystem
}

func (fs internalFs) Open(filename string) (ReadonlyFile, error) {
	return fs.Filesystem.Open(filename)
}

func (fs internalFs) Close() (_ error) {
	return
}
