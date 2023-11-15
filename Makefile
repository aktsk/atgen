NAME := atgen
VERSION = $(shell gobump show -r)
REVISION := $(shell git rev-parse --short HEAD)

all: build

setup:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/tcnksm/ghr
	go get github.com/Songmu/goxz/cmd/goxz
	go get github.com/x-motemen/gobump/cmd/gobump
	go get github.com/Songmu/ghch/cmd/ghch

test:
	go test ./lib
	go test -race ./lib

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
        -d=./pkg/dist/v${VERSION} \
        -n ${NAME}

release: package
	ghr -u aktsk v${VERSION} ./pkg/dist/v${VERSION}

bump:
	@sh -c "'$(CURDIR)/scripts/bump.sh'"
