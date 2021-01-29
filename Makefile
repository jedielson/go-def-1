GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
LDFLAGS += -s -w
SRCDIR ?= .
COMMANDS=$(wildcard ${SRCDIR}/cmd/*)
COMMANDS_BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))

GITFLAGS ?= GIT_DIR=${SRCDIR}/.git GIT_WORK_TREE=${SRCDIR}
ifeq ($(NOGIT),1)
  GIT_SUMMARY ?= Unknown
  GIT_BRANCH ?= Unknown
  GIT_MERGE ?= Unknown
else
  GIT_SUMMARY ?= $(shell ${GITFLAGS} git describe --tags --dirty --always)
  GIT_BRANCH ?= $(shell ${GITFLAGS} git symbolic-ref -q --short HEAD)
  GIT_MERGE ?= $(shell ${GITFLAGS} git rev-list --count --merges main)
endif

LDFLAGS += -X main.GitBranch=${GIT_BRANCH} -X main.GitSummary=${GIT_SUMMARY} -X main.GitMerge=${GIT_MERGE}

default: help 

## clean: cleans bin folder 
.PHONY: clean
clean: 
	@rm -rf bin/*
	
## deps: downloads mod dependencies
.PHONY: deps 
deps: 
	@go mod download

## build: builds cmd files 
.PHONY: build
build: build-cmd

## build-cmd: builds only cmd main files
.PHONY: build-cmd 
build-cmd:
	@echo GOOS       : $(GOOS)
	@echo GOARCH     : $(GOARCH)
	@echo LDFLAGS    : $(LDFLAGS)
	@echo CGO_ENABLED: $(CGO_ENABLED)
	@echo GIT_SUMMARY: $(GIT_SUMMARY)
	@echo GIT_BRANCH : $(GIT_BRANCH)
	@echo GIT_MERGE  : $(GIT_MERGE)
	@for dir in `ls cmd`; do \
      echo building: $$dir; \
			CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="${LDFLAGS}" -o bin/$$dir ./cmd/$$dir; \
	done

## build-docker: builds dockerfiles
.PHONY: build-docker
build-docker:
	@echo GOOS       : $(GOOS)
	@echo GOARCH     : $(GOARCH)
	@echo LDFLAGS    : $(LDFLAGS)
	@echo CGO_ENABLED: $(CGO_ENABLED)
	@echo GIT_SUMMARY: $(GIT_SUMMARY)
	@echo GIT_BRANCH : $(GIT_BRANCH)
	@echo GIT_MERGE  : $(GIT_MERGE)
	@for dir in `ls build`; do \
		echo dockering: $$dir; \
		docker build \
			--build-arg \
			-t jedielson/bookstore-$$dir:latest \
			-t jedielson/bookstore-$$dir:$(GIT_SUMMARY) \
			-f ./build/$$dir/Dockerfile .; \
	done

## test: test all files recursively 
.PHONY: test
test:
	@go test -v ./... 

## test-unit: run all unit tests
.PHONY: test-unit
test-unit:
	@go test -v -run Unit ./...


## test-integration: run all integration tests
.PHONY: test-integration
test-integration:
	@go test -v -run Integration ./...


## vet: vet all files recursively
vet:
	@go vet -v ./...

## all: runs clean test build 
.PHONY: all
all: clean deps test build

## help: show this help
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo