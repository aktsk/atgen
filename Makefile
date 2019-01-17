NAME := atgen
VERSION = $(shell gobump show -r)
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
	go build -o ./bin/$(NAME)

clean:
	rm bin/$(NAME)

package: setup
	@sh -c "'$(CURDIR)/scripts/package.sh'"

crossbuild: setup
	goxz -pv=v${VERSION} -build-ldflags="-X main.GitCommit=${REVISION}" \
        -arch=386,amd64 -d=./pkg/dist/v${VERSION} \
        -n ${NAME}

release: package
	ghr -u aktsk v${VERSION} ./pkg/dist/v${VERSION}

bump:
	@sh -c "'$(CURDIR)/scripts/bump.sh'"
