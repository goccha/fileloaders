module github.com/goccha/fileloaders/github-loader

go 1.21

require (
	github.com/goccha/fileloaders v0.0.1-alpha.7
	github.com/google/go-github/v66 v66.0.0
)

require github.com/google/go-querystring v1.1.0 // indirect

replace github.com/goccha/fileloaders => ./..
