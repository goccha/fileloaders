module test

go 1.21.3

require (
	cloud.google.com/go/storage v1.42.0
	github.com/aws/aws-sdk-go-v2 v1.27.2
	github.com/aws/aws-sdk-go-v2/config v1.26.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.55.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.50.6
	github.com/aws/smithy-go v1.20.2
	github.com/goccha/fileloaders v0.0.1-alpha.4
	github.com/goccha/fileloaders/github-loader v0.0.0-00010101000000-000000000000
	github.com/goccha/fileloaders/gs-loader v0.0.0-20200522141810-8b9b9c9b1b0e
	github.com/goccha/fileloaders/http-loader v0.0.0-20200522141810-8b9b9c9b1b0e
	github.com/goccha/fileloaders/s3-loader v0.0.0-20200522141810-8b9b9c9b1b0e
	github.com/goccha/fileloaders/ssm-loader v0.0.0-20200522141810-8b9b9c9b1b0e
	github.com/google/go-github/v62 v62.0.0
	google.golang.org/api v0.183.0
)

require (
	cloud.google.com/go v0.114.0 // indirect
	cloud.google.com/go/auth v0.5.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.8 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.12 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.52.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.52.0 // indirect
	go.opentelemetry.io/otel v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20240610135401-a8a62080eff3 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240610135401-a8a62080eff3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace (
	github.com/goccha/fileloaders => ../
	github.com/goccha/fileloaders/github-loader => ../github-loader
	github.com/goccha/fileloaders/gs-loader => ../gs-loader
	github.com/goccha/fileloaders/http-loader => ../http-loader
	github.com/goccha/fileloaders/s3-loader => ../s3-loader
	github.com/goccha/fileloaders/ssm-loader => ../ssm-loader
)
