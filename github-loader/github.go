package githubloader

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccha/fileloaders"
	"github.com/google/go-github/v62/github"
)

type Loader struct {
	client *github.Client
}

type LoaderBuilder struct {
	client *github.Client
}

func (b *LoaderBuilder) Load(ctx context.Context, path string, opt ...fileloaders.LoaderOption) (*fileloaders.File, error) {
	loader := &Loader{}
	for _, v := range opt {
		v(loader)
	}
	if loader.client == nil {
		loader.client = b.client
	}
	return loader.Load(ctx, path)
}
func (b *LoaderBuilder) List(ctx context.Context, path string, opt ...fileloaders.LoaderOption) ([]string, error) {
	loader := &Loader{}
	for _, v := range opt {
		v(loader)
	}
	if loader.client == nil {
		loader.client = b.client
	}
	return loader.List(ctx, path)
}

func WithAuthToken(token string) fileloaders.LoaderOption {
	return func(l fileloaders.Loader) {
		if v, ok := l.(*Loader); ok {
			v.client = github.NewClient(http.DefaultClient).WithAuthToken(token)
		}
	}
}

type Builder func() *github.Client

func Load(ctx context.Context, c *github.Client, path string) (*fileloaders.File, error) {
	file := fileloaders.Parse(path)
	if file == nil || file.Type != "github" || file.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	repoIndex := strings.Index(file.Path, "/")
	repo := file.Path[:repoIndex]
	filepath := file.Path[repoIndex+1:]
	var opts *github.RepositoryContentGetOptions
	if index := strings.LastIndex(filepath, "?"); index > 0 {
		if query, err := url.ParseQuery(filepath[index+1:]); err != nil {
			return nil, err
		} else if query.Has("ref") {
			opts = &github.RepositoryContentGetOptions{Ref: query.Get("ref")}
		}
		filepath = filepath[:index]
	}
	fileContent, _, res, err := c.Repositories.GetContents(ctx, file.Bucket, repo, filepath, opts)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	v, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}
	return file.Add(fileloaders.WithHash(fileContent.SHA)).WriteBody([]byte(v)), nil
}

func List(ctx context.Context, c *github.Client, path string) ([]string, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "github" || filePath.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	repoIndex := strings.Index(filePath.Path, "/")
	repo := filePath.Path[:repoIndex]
	shaIndex := strings.Index(filePath.Path[repoIndex+1:], "/")
	sha := ""
	if shaIndex < 0 {
		sha = filePath.Path[repoIndex+1:]
	} else {
		sha = filePath.Path[repoIndex+1 : shaIndex]
	}

	tree, res, err := c.Git.GetTree(ctx, filePath.Bucket, repo, sha, false)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	result := make([]string, len(tree.Entries))
	for i, v := range tree.Entries {
		result[i] = *v.Path
	}
	return result, nil
}

func New(c *github.Client) fileloaders.Loader {
	return &LoaderBuilder{
		client: c,
	}
}

func (l *Loader) Load(ctx context.Context, path string, opt ...fileloaders.LoaderOption) (*fileloaders.File, error) {
	return Load(ctx, l.client, path)
}

func (l *Loader) List(ctx context.Context, path string, opt ...fileloaders.LoaderOption) ([]string, error) {
	return List(ctx, l.client, path)
}

func WithClient(api *github.Client) fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		m["github"] = New(api)
	}
}

func With() fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		m["github"] = New(nil)
	}
}
