package s3loader

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/goccha/fileloaders"
)

type Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func Load(ctx context.Context, api Client, path string) (*fileloaders.File, error) {
	file := fileloaders.Parse(path)
	if file == nil || file.Type != "s3" || file.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	var version *string
	path = file.Path
	if index := strings.LastIndex(path, "?"); index > 0 {
		if query, err := url.ParseQuery(path[index+1:]); err != nil {
			return nil, err
		} else if query.Has("version") {
			version = aws.String(query.Get("version"))
		}
		path = path[:index]
	}
	result, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket:    aws.String(file.Bucket),
		Key:       aws.String(path),
		VersionId: version,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = result.Body.Close()
	}()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return file.WriteBody(body).Add(
		fileloaders.WithHash(result.ETag),
		fileloaders.WithVersion(result.VersionId)), nil
}

func List(ctx context.Context, api Client, path string) ([]string, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "s3" || filePath.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	in := &s3.ListObjectsV2Input{
		Bucket: aws.String(filePath.Bucket),
	}
	if filePath.Path != "" {
		in.Prefix = aws.String(filePath.Path)
	}
	out, err := api.ListObjectsV2(ctx, in)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(out.Contents))
	for i, v := range out.Contents {
		result[i] = *v.Key
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
	return &Loader{client: api}
}

func With(api Client) fileloaders.Option {
	return func(m map[string]fileloaders.Loader) {
		m["s3"] = New(api)
	}
}
