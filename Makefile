NAME := atge
VERSION = $(shell gobump show -r ./version)
REVISION := $(shell git rev-parse --short HEAD)

all: build

setup:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/tcnksm/ghr
	go get github.com/Songmu/goxz/cmd/goxz
	go get github.com/motemen/gobump/cmd/gobump

test: lint
	go test ./lib
	go test -race ./lib

lint: setup
	golint ./...

fmt: setup
	goimports -w .

build:
	go build
