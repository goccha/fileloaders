module github.com/goccha/fileloaders/s3

go 1.21.3

require (
	github.com/aws/aws-sdk-go-v2 v1.23.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.43.0
	github.com/goccha/fileloaders v0.0.0-20231021043511-2afe32ef1c70
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.3 // indirect
	github.com/aws/smithy-go v1.17.0 // indirect
)

replace github.com/goccha/fileloaders => ./..
