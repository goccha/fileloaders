package httploader

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/goccha/fileloaders"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

func Load(c Client, path string) ([]byte, error) {
	if c == nil {
		c = http.DefaultClient
	}
	res, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

type Loader struct {
	client Client
}

func (l *Loader) Load(ctx context.Context, path string) ([]byte, error) {
	return Load(l.client, path)
}
func (l *Loader) List(ctx context.Context, path string) ([]string, error) {
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
