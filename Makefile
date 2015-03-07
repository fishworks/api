SHELL = /bin/bash

GO = go
INSTALL = install
RM = rm -f

GOBUILD = $(GO) build -o
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test -cover --race
GOVET = $(GO) vet

prefix := /usr/local

# The directory to install in
bin_dir = $(prefix)/bin

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/bacongobbler/api

# used to reference the output directory for build artifacts
build_dir = bin

SRC_PACKAGES = api server
REPO_SRC_PACKAGES = $(addprefix $(repo_path)/,$(SRC_PACKAGES))

all: build

build:
	$(GOBUILD) $(build_dir)/api $(repo_path)/api

clean:
	$(RM) $(build_dir)/*

install:
	$(INSTALL) -c $(build_dir)/api $(bin_dir)/api

test:
# display output, then check
	$(GOFMT) $(SRC_PACKAGES)
	@$(GOFMT) $(SRC_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi

# display output, then check
	$(GOLINT) ./...
	$(GOLINT) ./... | read; if [ $$? == 0 ]; then echo "golint check failed."; exit 1; fi

	$(GOTEST) $(REPO_SRC_PACKAGES)
	$(GOVET) $(REPO_SRC_PACKAGES)

.PHONY: all build clean install test
