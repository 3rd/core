SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

dev: ## watch for changes and run tests
	watchexec --restart -c reset --exts go -- "go test ./..."
.PHONY: dev

dev-task-interactive: ## dev: task interactive
	watchexec --stop-timeout=0 --no-process-group --stop-signal SIGKILL -e go -r -c clear -w . -- "go run . task interactive"
.PHONY: dev

test: ## run tests
	go test ./...
.PHONY: test

build: ## build
	go build .
.PHONY: build

install: ## install
	go install .
.PHONY: install

bench: ## run benchmarks
	go test -bench=./... -benchmem
.PHONY: bench

cover: ## run tests and show coverage
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: cover

cpu: ## pprof
	go test -cpuprofile cpu.prof -bench ./...
	go tool pprof -http=":8081" cpu.prof
.PHONY: cpu

mem: ## pprof
	go test -memprofile mem.prof -bench ./...
	go tool pprof -http=":8081" mem.prof
.PHONY: mem
