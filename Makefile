###############################################################################
# VARIABLES
###############################################################################

export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

PKG_PATH=github.com/nikoksr/proji
GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_REV=$(shell git rev-parse --short HEAD)

###############################################################################
# DEPENDENCIES
###############################################################################

setup:
	go mod tidy
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: setup

###############################################################################
# TESTS
###############################################################################

test:
	go test -failfast -race ./...
.PHONY: test

gen-coverage:
	@go test -race -covermode=atomic -coverprofile=coverage.out ./... > /dev/null
.PHONY: gen-coverage

coverage: gen-coverage
	go tool cover -func coverage.out
.PHONY: coverage

coverage-html: gen-coverage
	go tool cover -html=coverage.out -o cover.html
.PHONY: coverage-html



###############################################################################
# CODE HEALTH
###############################################################################

fmt:
	@gofumpt -w -l . > /dev/null

	@goimports -w -l -local github.com/nikoksr/proji . > /dev/null

	@gci write --section standard --section default --section "Prefix(github.com/nikoksr/proji)" . > /dev/null
.PHONY: fmt

lint:
	@golangci-lint run --config .golangci.yml
.PHONY: lint

ci: fmt test lint
.PHONY: ci

###############################################################################
# BUILDS
###############################################################################

prepare-build:
	mkdir -p ./bin/debug ./bin/release ./docs
.PHONY: prepare-build

check-optimizations: prepare-build
	go build -gcflags='-m -m' -o ./bin/debug/proji ./cmd/proji
.PHONY: check-optimizations

build-debug: prepare-build
	CGO_ENABLED=0 go build -o ./bin/debug/proji ./cmd/proji
.PHONY: build-debug

build-debug-win: prepare-build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/debug/proji-win-amd64.exe ./cmd/proji
.PHONY: build-debug-win

build-release: prepare-build
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/release/proji ./cmd/proji
.PHONY: build-release

install:
	CGO_ENABLED=0 go install ./cmd/proji
.PHONY: install

clean:
	rm -rf ./bin 2> /dev/null
	rm cover.html coverage.out 2> /dev/null
.PHONY: clean

.DEFAULT_GOAL := install
