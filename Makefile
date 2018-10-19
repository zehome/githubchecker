all: build

get:
	go get -u "github.com/hashicorp/go-version"
	go get -u "github.com/tcnksm/go-latest"


build:
	go build -o githubchecker main.go
