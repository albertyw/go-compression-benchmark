.PHONY: all build test benchmark vet

all: build

build:
	go build -o benchmark .

test:
	go test ./...

vet:
	go vet ./...

benchmark: build
	./benchmark
