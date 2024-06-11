package s3loader

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/goccha/fileloaders"
)

type Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func Load(ctx context.Context, api Client, path string) ([]byte, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "s3" || filePath.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	result, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(filePath.Bucket),
		Key:    aws.String(filePath.Path),
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = result.Body.Close()
	}()
	return io.ReadAll(result.Body)
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

func (l *Loader) Load(ctx context.Context, path string, opt ...fileloaders.LoaderOption) ([]byte, error) {
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
