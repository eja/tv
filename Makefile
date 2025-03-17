.PHONY: clean test lint tv modulator 

GOLANG_CROSS_VERSION := v1.20
GOPATH ?= '$(HOME)/go'

all: lint tv modulator

clean:
	@rm -f tv tv.exe
	@rm -f modulator modulator.exe

lint:
	@gofmt -w .

test:
	@go mod tidy
	@go mod verify
	@go vet ./...
	@go test -v ./test

tv:
	@go build -ldflags "-s -w" -o tv ./cmd/tv/...
	@strip tv

modulator:
	@go build -ldflags "-s -w" -o modulator ./cmd/modulator/...
	@strip modulator
