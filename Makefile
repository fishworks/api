SHELL = /bin/bash
GO15VENDOREXPERIMENT = 1

GO = go
INSTALL = install
RM = rm -f

GOBUILD = $(GO) build -o
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test --cover --race
GOVET = $(GO) vet

prefix := /usr/local

# The directory to install in
bin_dir = $(prefix)/bin

# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/fishworks/api

# used to reference the output directory for build artifacts
build_dir = bin

NOVENDOR = $(shell glide nv)

all: build

build:
	$(GOBUILD) $(build_dir)/api $(repo_path)/cmd

clean:
	$(RM) $(build_dir)/*

install:
	$(INSTALL) -c $(build_dir)/api $(bin_dir)/api

test:
	$(GOTEST) $(NOVENDOR)

.PHONY: all build clean install test
