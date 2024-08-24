package httploader

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/goccha/fileloaders"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

func Load(c Client, path string) (*fileloaders.File, error) {
	if c == nil {
		c = http.DefaultClient
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	res, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	file := &fileloaders.File{
		Type:   u.Scheme,
		Bucket: u.Host,
		Path:   u.Path,
	}
	return file.WriteBody(body), nil
}

type Loader struct {
	client Client
}

func (l *Loader) Load(ctx context.Context, path string, opt ...fileloaders.LoaderOption) (*fileloaders.File, error) {
	return Load(l.client, path)
}
func (l *Loader) List(ctx context.Context, path string, opt ...fileloaders.LoaderOption) ([]string, error) {
	return nil, fileloaders.ErrNotSupported
}

func New(c Client) *Loader {
	return &Loader{
		client: c,
	}
}

func With(c Client) fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		loader := New(c)
		m["http"] = loader
		m["https"] = loader
	}
}
