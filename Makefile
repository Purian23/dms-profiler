.PHONY: fmt test tidy all

fmt:
	go fmt ./...

test:
	go test ./...

tidy:
	go mod tidy

all: fmt test
