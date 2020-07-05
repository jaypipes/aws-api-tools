GO111MODULE=on

.PHONY: all test

all: test

test:
	go test -tags codegen -v ./...

build:
	# -tags codegen is required due to importing aws-sdk-go/private packages
	go build -tags codegen -o bin/aws-api-tool cmd/aws-api-tool/main.go
