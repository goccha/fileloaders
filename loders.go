package fileloaders

import (
	"context"
	"errors"
	"os"
	"strings"
)

var ErrNotSupported = errors.New("does not supported")

type LoaderFunc func(ctx context.Context, path string) (*File, error)
type ListFunc func(ctx context.Context, path string) ([]string, error)

type Option func(m map[string]Loader)

func Setup(options ...Option) {
	if root == nil {
		root = New(options...)
	} else {
		for _, option := range options {
			option(root.loaders)
		}
	}
}

func New(options ...Option) *MapLoader {
	loader := &MapLoader{
		loaders: make(map[string]Loader),
	}
	for _, option := range options {
		option(loader.loaders)
	}
	return loader
}

var root *MapLoader

func Load(ctx context.Context, path string, opt ...LoaderOption) (*File, error) {
	if root != nil {
		if v, err := root.Load(ctx, path, opt...); err != nil {
			if !errors.Is(err, ErrNotSupported) {
				return nil, err
			}
		} else {
			return v, nil
		}
	}
	return LoadFile(ctx, path)
}

func List(ctx context.Context, path string, opt ...LoaderOption) ([]string, error) {
	if root != nil {
		if v, err := root.List(ctx, path, opt...); err != nil {
			if !errors.Is(err, ErrNotSupported) {
				return nil, err
			}
		} else {
			return v, nil
		}
	}
	return ListFile(ctx, path)
}

func LoadFile(ctx context.Context, path string) (*File, error) {
	path = strings.TrimPrefix(path, "file://")
	bin, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &File{
		Type: "file",
		Path: path,
		body: bin,
	}, nil
}

func ListFile(ctx context.Context, path string) ([]string, error) {
	path = strings.TrimPrefix(path, "file://")
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(entries))
	for i, v := range entries {
		result[i] = v.Name()
	}
	return result, nil
}

type LoaderOption func(loader Loader)

type Loader interface {
	Load(ctx context.Context, path string, opt ...LoaderOption) (*File, error)
	List(ctx context.Context, path string, opt ...LoaderOption) ([]string, error)
}

type MapLoader struct {
	loaders map[string]Loader
}

func (m *MapLoader) Load(ctx context.Context, path string, opt ...LoaderOption) (*File, error) {
	index := strings.Index(path, "://")
	var prefix string
	if index > 0 {
		prefix = path[:index]
	}
	loader, ok := m.loaders[prefix]
	if !ok {
		if prefix == "file" {
			return LoadFile(ctx, path)
		}
		return nil, ErrNotSupported
	}
	return loader.Load(ctx, path, opt...)
}

func (m *MapLoader) List(ctx context.Context, path string, opt ...LoaderOption) ([]string, error) {
	index := strings.Index(path, "://")
	var prefix string
	if index > 0 {
		prefix = path[:index]
	}
	loader, ok := m.loaders[prefix]
	if !ok {
		if prefix == "file" {
			return ListFile(ctx, path)
		}
		return nil, ErrNotSupported
	}
	return loader.List(ctx, path, opt...)
}
