NAME    := wikiracer
VERSION := v0.1.0
BUILD   := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X main.version=${VERSION} -X main.build=${BUILD}"

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin
	rm -rf vendor
	rm -rf log

.PHONY: setup
setup:
	go get -u github.com/golang/dep/cmd/dep

.PHONY: deps
deps: setup
	dep ensure

.PHONY: install
install: deps
	go install $(LDFLAGS)
