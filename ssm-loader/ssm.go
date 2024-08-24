package ssmloader

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/goccha/fileloaders"
)

type Client interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
}

func Load(ctx context.Context, api Client, path string) (*fileloaders.File, error) {
	file := fileloaders.Parse(path)
	if file == nil || file.Type != "ssm" || file.Bucket == "" {
		return nil, fileloaders.ErrNotSupported
	}
	var key string
	if file.Bucket != "" {
		key = "/" + file.Bucket + "/" + file.Path
	} else {
		key = "/" + file.Path
	}
	out, err := api.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if out.Parameter != nil {
		if out.Parameter.Version > 0 {
			version := strconv.Itoa(int(out.Parameter.Version))
			file = file.Add(fileloaders.WithVersion(&version))
		}
		if out.Parameter.Value != nil {
			return file.WriteBody([]byte(*out.Parameter.Value)), nil
		}
	}
	return &fileloaders.File{}, nil
}

func List(ctx context.Context, api Client, path string) ([]string, error) {
	filePath := fileloaders.Parse(path)
	if filePath == nil || filePath.Type != "ssm" {
		return nil, fileloaders.ErrNotSupported
	}
	var key string
	if filePath.Bucket != "" {
		key = "/" + filePath.Bucket + "/"
		if filePath.Path != "" {
			key = key + filePath.Path
		}
	} else {
		key = "/" + filePath.Path
	}
	input := &ssm.DescribeParametersInput{}
	if key != "" {
		input = &ssm.DescribeParametersInput{
			ParameterFilters: []types.ParameterStringFilter{
				{
					Key:    aws.String(string(types.ParametersFilterKeyName)),
					Option: aws.String("BeginsWith"),
					Values: []string{key},
				},
			},
		}
	}
	out, err := api.DescribeParameters(ctx, input)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(out.Parameters))
	for _, v := range out.Parameters {
		result = append(result, *v.Name)
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
		m["ssm"] = New(api)
	}
}
