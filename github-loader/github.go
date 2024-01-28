package githubloader

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/goccha/fileloaders"
	"github.com/google/go-github/v56/github"
)

type Loader struct {
	client *github.Client
}

func Load(ctx context.Context, c *github.Client, path string) ([]byte, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "github" || filePath.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	repoIndex := strings.Index(filePath.Path, "/")
	repo := filePath.Path[:repoIndex]
	filepath := filePath.Path[repoIndex+1:]

	file, _, res, err := c.Repositories.GetContents(ctx, filePath.Bucket, repo, filepath, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	//defer file.Close()
	//return io.ReadAll(file)
	v, err := file.GetContent()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
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

func New(c *github.Client) *Loader {
	return &Loader{
		client: c,
	}
}

func (l *Loader) Load(ctx context.Context, path string) ([]byte, error) {
	return Load(ctx, l.client, path)
}

func (l *Loader) List(ctx context.Context, path string) ([]string, error) {
	return List(ctx, l.client, path)
}

func With(api *github.Client) fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		m["github"] = New(api)
	}
}
