Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
SOURCE_DIRS = cmd pkg main.go
export GO111MODULE=on

.PHONY: all
all: gofmt test build dist

.PHONY: build
build:
	go build

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make gofmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  -a -installsuffix cgo -o bin/rekind
	CGO_ENABLED=0 GOOS=darwin go build -ldflags  -a -installsuffix cgo -o bin/rekind-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags  -installsuffix cgo -o bin/rekind-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags  -a -installsuffix cgo -o bin/rekind-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags  -a -installsuffix cgo -o bin/rekind-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags  -a -installsuffix cgo -o bin/rekind.exe
