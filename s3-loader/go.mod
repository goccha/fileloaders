module github.com/goccha/fileloaders/s3-loader

go 1.21.3

require (
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.50.0
	github.com/goccha/fileloaders v0.0.1-alpha.1
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.0 // indirect
	github.com/aws/smithy-go v1.20.0 // indirect
)

replace github.com/goccha/fileloaders => ./..
