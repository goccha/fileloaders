package testdata

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/goccha/fileloaders"
	"github.com/goccha/fileloaders/gs"
	httploader "github.com/goccha/fileloaders/http"
	s3loader "github.com/goccha/fileloaders/s3"
	"google.golang.org/api/option"
)

func setupS3(ctx context.Context) error {
	region := "ap-northeast-1"
	//logLevel := aws.LogSigning | aws.LogRequestWithBody | aws.LogResponseWithBody | aws.LogRetries
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	var err error
	var cfg aws.Config
	if cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region),
		//config.WithClientLogMode(logLevel),
		config.WithLogger(logging.NewStandardLogger(os.Stdout)), config.WithLogConfigurationWarnings(true),
	); err != nil {
		return err
	}
	endpoint := "http://localhost:9090"
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	fileloaders.Setup(s3loader.With(client))
	if _, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("README.md"),
		Body:   strings.NewReader("# README"),
	}); err != nil {
		return err
	}
	return nil
}

func setupGs(ctx context.Context) (func(ctx context.Context) error, error) {
	_ = os.Setenv("STORAGE_EMULATOR_HOST", "localhost:8000")
	var err error
	client, err := storage.NewClient(ctx, option.WithEndpoint("http://localhost:8000/storage/v1/"), option.WithoutAuthentication())
	if err != nil {
		panic(err)
	}
	fileloaders.Setup(gs.With(client))
	projectID := "test-project"
	bucketName := "test-bucket"
	if err := client.Bucket(bucketName).Create(ctx, projectID, nil); err != nil {
		return nil, err
	}
	w := client.Bucket(bucketName).Object("README.md").NewWriter(ctx)
	w.ContentType = "text/markdown"
	_, err = io.Copy(w, bytes.NewReader([]byte("# README")))
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return func(ctx context.Context) error {
		if err := client.Bucket(bucketName).Object("README.md").Delete(ctx); err != nil {
			return err
		}
		if err := client.Bucket(bucketName).Delete(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	if err := setupS3(ctx); err != nil {
		t.Fatal(err)
	}
	if clean, err := setupGs(ctx); err != nil {
		t.Fatal(err)
	} else {
		defer func() {
			if err := clean(ctx); err != nil {
				t.Fatal(err)
			}
		}()
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/README.md" {
			_, _ = fmt.Fprint(w, "# README")
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	// 別goroutine上でリッスンが開始される
	ts := httptest.NewServer(h)
	defer ts.Close()
	fileloaders.Setup(httploader.With(http.DefaultClient))

	list, err := fileloaders.List(ctx, "s3://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatal("invalid s3 list")
	}
	list, err = fileloaders.List(ctx, "gs://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatal("invalid gs list")
	}

	bin, err := fileloaders.Load(context.Background(), "s3://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
	bin, err = fileloaders.Load(context.Background(), "gs://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
	bin, err = fileloaders.Load(context.Background(), ts.URL+"/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
	bin, err = fileloaders.Load(context.Background(), "../README.md")
	if err != nil {
		t.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(bin))
	if s.Scan() {
		v := s.Text()
		if v != "# fileloaders" {
			t.Fatal("invalid content")
		}
	} else {
		t.Fatal("invalid load")
	}
}
