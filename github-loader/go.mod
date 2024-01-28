module github.com/goccha/fileloaders/github-loader

go 1.21

require (
	github.com/goccha/fileloaders v0.0.0-20231118034152-b01698c5e160
	github.com/google/go-github/v56 v56.0.0
)

require github.com/google/go-querystring v1.1.0 // indirect

replace github.com/goccha/fileloaders => ./..
