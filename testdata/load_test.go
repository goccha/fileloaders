package testdata

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go/logging"
	"github.com/goccha/fileloaders"
	"github.com/goccha/fileloaders/github-loader"
	"github.com/goccha/fileloaders/gs-loader"
	"github.com/goccha/fileloaders/http-loader"
	"github.com/goccha/fileloaders/s3-loader"
	"github.com/goccha/fileloaders/ssm-loader"
	"github.com/google/go-github/v56/github"
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

func setupSsm(ctx context.Context) error {
	region := "ap-northeast-1"
	//logLevel := aws.LogSigning | aws.LogRequestWithBody | aws.LogResponseWithBody | aws.LogRetries
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	var err error
	var cfg aws.Config
	if cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region),
		config.WithLogger(logging.NewStandardLogger(os.Stdout)), config.WithLogConfigurationWarnings(true),
	); err != nil {
		return err
	}
	endpoint := "http://localhost:4566"
	client := ssm.NewFromConfig(cfg, func(o *ssm.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	fileloaders.Setup(ssmloader.With(client))
	return nil
}

func setupGs(ctx context.Context) (func(ctx context.Context) error, error) {
	_ = os.Setenv("STORAGE_EMULATOR_HOST", "localhost:8000")
	var err error
	client, err := storage.NewClient(ctx, option.WithEndpoint("http://localhost:8000/storage/v1/"), option.WithoutAuthentication())
	if err != nil {
		panic(err)
	}
	fileloaders.Setup(gsloader.With(client))
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

func setupGithub(ctx context.Context, baseUrl *url.URL) error {
	token := os.Getenv("GITHUB_TOKEN")
	cli := github.NewClient(http.DefaultClient).WithAuthToken(token)
	var err error
	if token == "" {
		if v := os.Getenv("GITHUB_TEST_HOST"); v != "" {
			cli.BaseURL, err = url.Parse(v)
			if err != nil {
				return err
			}
		} else if baseUrl != nil {
			cli.BaseURL = baseUrl
		}
	}
	fileloaders.Setup(githubloader.With(cli))
	return nil
}

func TestLoad(t *testing.T) {
	bin, err := fileloaders.Load(context.Background(), "../README.md")
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

func TestHttp(t *testing.T) {
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

	bin, err := fileloaders.Load(context.Background(), ts.URL+"/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
}

func TestS3(t *testing.T) {
	ctx := context.Background()
	if err := setupS3(ctx); err != nil {
		t.Fatal(err)
	}
	list, err := fileloaders.List(ctx, "s3://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatal("invalid s3 list")
	}

	bin, err := fileloaders.Load(context.Background(), "s3://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
}

func TestGs(t *testing.T) {
	ctx := context.Background()
	if clean, err := setupGs(ctx); err != nil {
		t.Fatal(err)
	} else {
		defer func() {
			if err := clean(ctx); err != nil {
				t.Fatal(err)
			}
		}()
	}
	list, err := fileloaders.List(ctx, "gs://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatal("invalid gs list")
	}
	bin, err := fileloaders.Load(context.Background(), "gs://test-bucket/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# README" {
		t.Fatal("invalid load")
	}
}

func TestSsm(t *testing.T) {
	ctx := context.Background()
	if err := setupSsm(ctx); err != nil {
		t.Fatal(err)
	}
	list, err := fileloaders.List(ctx, "ssm://parameter/")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatal("invalid ssm list")
	}
	bin, err := fileloaders.Load(context.Background(), "ssm://parameter1/test")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "value/fileloaders/test" {
		t.Fatal("invalid load")
	}
	bin, err = fileloaders.Load(context.Background(), "ssm://parameter1/secure")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "secure-value/fileloaders/test" {
		t.Fatal("invalid load")
	}

}

func TestGithub(t *testing.T) {
	ctx := context.Background()
	baseUrl := ""
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/goccha/fileloaders/git/trees/main" {
			_, _ = fmt.Fprint(w, `{
	"sha": "b01698c5e1607052a0be6439b355625e37c7e769",
	"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/b01698c5e1607052a0be6439b355625e37c7e769",
	"tree": [
		{
			"path": ".github",
			"mode": "040000",
			"type": "tree",
			"sha": "5d8ca12ef91b5696511d36ea2178b4c076daaab2",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/5d8ca12ef91b5696511d36ea2178b4c076daaab2"
		},
		{
			"path": ".gitignore",
			"mode": "100644",
			"type": "blob",
			"sha": "723ef36f4e4f32c4560383aa5987c575a30c6535",
			"size": 5,
			"url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/723ef36f4e4f32c4560383aa5987c575a30c6535"
		},
		{
			"path": "LICENSE",
			"mode": "100644",
			"type": "blob",
			"sha": "9afb442af50eebae8768d568d44e790da8771695",
			"size": 1084,
			"url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/9afb442af50eebae8768d568d44e790da8771695"
		},
		{
			"path": "README.md",
			"mode": "100644",
			"type": "blob",
			"sha": "b6c811e48fc707be5a6e01ab3ae50b361990b5a7",
			"size": 13,
			"url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/b6c811e48fc707be5a6e01ab3ae50b361990b5a7"
		},
		{
			"path": "docker",
			"mode": "040000",
			"type": "tree",
			"sha": "d62173472d5b89d63fd93d2e6f6f19e6bb76f008",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/d62173472d5b89d63fd93d2e6f6f19e6bb76f008"
		},
		{
			"path": "go.mod",
			"mode": "100644",
			"type": "blob",
			"sha": "a1427fc3e764c6f9150e39a204af5f23bbecfd9a",
			"size": 46,
			"url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/a1427fc3e764c6f9150e39a204af5f23bbecfd9a"
		},
		{
			"path": "gs",
			"mode": "040000",
			"type": "tree",
			"sha": "9bd61488603fa965bf7fee96b32df88ed80b2b80",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/9bd61488603fa965bf7fee96b32df88ed80b2b80"
		},
		{
			"path": "http",
			"mode": "040000",
			"type": "tree",
			"sha": "eb56fb5aa19e82785faafd69554c992e3fcd7fe4",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/eb56fb5aa19e82785faafd69554c992e3fcd7fe4"
		},
		{
			"path": "loders.go",
			"mode": "100644",
			"type": "blob",
			"sha": "01941774ddbde07afc862a6128e88cd3e3401ea8",
			"size": 2518,
			"url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/01941774ddbde07afc862a6128e88cd3e3401ea8"
		},
		{
			"path": "s3",
			"mode": "040000",
			"type": "tree",
			"sha": "fa0bcbfec8feadf927e1bad597d059073af203a7",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/fa0bcbfec8feadf927e1bad597d059073af203a7"
		},
		{
			"path": "testdata",
			"mode": "040000",
			"type": "tree",
			"sha": "8eadd6edd6d53eb7649b923933e983e6a44db180",
			"url": "https://api.github.com/repos/goccha/fileloaders/git/trees/8eadd6edd6d53eb7649b923933e983e6a44db180"
		}
	],
	"truncated": false
}`)
		} else if strings.HasPrefix(r.URL.Path, "/repos/goccha/fileloaders/contents") {
			body := `{
	"name": "README.md",
	"path": "README.md",
	"sha": "b6c811e48fc707be5a6e01ab3ae50b361990b5a7",
	"size": 13,
	"url": "https://api.github.com/repos/goccha/fileloaders/contents/README.md?ref=main",
	"html_url": "https://github.com/goccha/fileloaders/blob/main/README.md",
	"git_url": "https://api.github.com/repos/goccha/fileloaders/git/blobs/b6c811e48fc707be5a6e01ab3ae50b361990b5a7",
	"download_url": "https://raw.githubusercontent.com/goccha/fileloaders/main/README.md",
	"type": "file",
	"content": "IyBmaWxlbG9hZGVycw==\n",
	"encoding": "base64",
	"_links": {
		"self": "https://api.github.com/repos/goccha/fileloaders/contents/README.md?ref=main",
		"git": "https://api.github.com/repos/goccha/fileloaders/git/blobs/b6c811e48fc707be5a6e01ab3ae50b361990b5a7",
		"html": "https://github.com/goccha/fileloaders/blob/main/README.md"
	}
}`
			//body = strings.ReplaceAll(body, "${BASE_URL}", baseUrl)
			_, _ = fmt.Fprint(w, body)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	// 別goroutine上でリッスンが開始される
	ts := httptest.NewServer(h)
	defer ts.Close()
	baseUrl = ts.URL + "/"
	base, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatal(err)
	}
	if err := setupGithub(ctx, base); err != nil {
		t.Fatal(err)
	}
	list, err := fileloaders.List(ctx, "github://goccha/fileloaders/main")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Fatal("invalid gs list")
	}

	bin, err := fileloaders.Load(ctx, "github://goccha/fileloaders/README.md")
	if err != nil {
		t.Fatal(err)
	}
	if string(bin) != "# fileloaders" {
		t.Fatal("invalid load")
	}
}
