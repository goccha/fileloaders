package s3

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/goccha/fileloaders"
)

type Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func Load(ctx context.Context, api Client, path string) ([]byte, error) {
	path = strings.TrimPrefix(path, "s3://")
	path = strings.TrimPrefix(path, "/")
	bucketIndex := strings.Index(path, "/")
	bucket := path[:bucketIndex]
	file := path[bucketIndex+1:]
	result, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
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
	path = strings.TrimPrefix(path, "s3://")
	path = strings.TrimPrefix(path, "/")
	bucketIndex := strings.Index(path, "/")
	bucket := path[:bucketIndex]
	prefix := path[bucketIndex+1:]
	in := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if prefix != "" {
		in.Prefix = aws.String(prefix)
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

func (l *Loader) Load(ctx context.Context, path string) ([]byte, error) {
	return Load(ctx, l.client, path)
}

func (l *Loader) List(ctx context.Context, path string) ([]string, error) {
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
