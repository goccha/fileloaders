module github.com/goccha/fileloaders/ssm-loader

go 1.21.3

require (
	github.com/aws/aws-sdk-go-v2 v1.25.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.49.0
	github.com/goccha/fileloaders v0.0.1-alpha.2
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.1 // indirect
	github.com/aws/smithy-go v1.20.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/goccha/fileloaders => ./..
