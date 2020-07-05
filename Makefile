GO111MODULE=on

.PHONY: all test

all: test

test:
	go test -v ./...
