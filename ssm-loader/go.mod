module github.com/goccha/fileloaders/ssm-loader

go 1.21.3

require (
	github.com/aws/aws-sdk-go-v2 v1.27.2
	github.com/aws/aws-sdk-go-v2/service/ssm v1.50.6
	github.com/goccha/fileloaders v0.0.1-alpha.4
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.9 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/goccha/fileloaders => ./..
