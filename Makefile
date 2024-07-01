.PHONY: clean test lint tv release-dry-run release

PACKAGE_NAME := tv
GOLANG_CROSS_VERSION := v1.20
GOPATH ?= '$(HOME)/go'

all: lint tv

clean:
	@rm -f tv tv.exe

lint:
	@gofmt -w .

test:
	@go mod tidy
	@go mod verify
	@go vet ./...
	@go test -v ./test

tv:
	@go build -ldflags "-s -w" -o tv cmd/tv/main.go
	@strip tv

release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip-validate
