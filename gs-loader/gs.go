package gsloader

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
	"github.com/goccha/fileloaders"
	"google.golang.org/api/iterator"
)

type Client interface {
	Bucket(name string) *storage.BucketHandle
}

func Load(ctx context.Context, api Client, path string) (*fileloaders.File, error) {
	file := fileloaders.Parse(path)
	if file == nil || file.Type != "gs" || file.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	bucketHandle := api.Bucket(file.Bucket)
	obj := bucketHandle.Object(file.Path)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return file.WriteBody(body), nil
}
func List(ctx context.Context, api Client, path string) ([]string, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "gs" || filePath.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	bucketHandle := api.Bucket(filePath.Bucket)
	iter := bucketHandle.Objects(ctx, &storage.Query{
		Prefix: filePath.Path,
	})
	var result []string
	for {
		obj, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		result = append(result, obj.Name)
	}
	return result, nil
}

type Loader struct {
	client Client
}

func (l *Loader) Load(ctx context.Context, path string, opt ...fileloaders.LoaderOption) (*fileloaders.File, error) {
	return Load(ctx, l.client, path)
}
func (l *Loader) List(ctx context.Context, path string, opt ...fileloaders.LoaderOption) ([]string, error) {
	return List(ctx, l.client, path)
}

func New(api Client) *Loader {
	return &Loader{
		client: api,
	}
}

func With(api Client) fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		m["gs"] = New(api)
	}
}
