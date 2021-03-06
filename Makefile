.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)

build: clean
	mkdir -p build
	go build -o build/busboy -ldflags "-X main.build=${hash}" cmd/busboy/main.go

test:
	go test ./...

run:
	go run cmd/busboy/main.go

clean:
	rm -rf build

