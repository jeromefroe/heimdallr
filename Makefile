PROJECT := heimdallr
REGISTRY ?= quay.io/jeromefroe
IMAGE := $(REGISTRY)/$(PROJECT)

GIT_REF := $(shell git rev-parse --short HEAD)
VERSION ?= $(GIT_REF)

SRC_DIRS := ./cmd ./pkg
SRC_FILES := $(shell find $(SRC_DIRS) -name '*.go')
PKGS := $(shell go list ./cmd/... ./pkg/... | grep -v /pkg/apis/ | grep -v /pkg/client/)

# Tools
retool:
	which retool || go get github.com/twitchtv/retool

# Dependencies
dep-install: retool
	retool do dep ensure -vendor-only -v

dep-update: retool
	retool do dep ensure -update -v

# Code Generation
gen: retool
	retool do go generate ./...
	@(bash hack/update-codegen.sh)

# Linting
fmt:
	@test -z "$(shell gofmt -s -l -d -e $(SRC_DIRS) | tee /dev/stderr)"

vet:
	go vet ./...

megacheck: retool
	retool do megacheck $(PKGS)

misspell: retool
	retool do misspell $(PKGS)

unconvert: retool
	retool do unconvert -v $(PKGS)

ineffassign: retool
	retool do ineffassign $(SRC_FILES)

unparam: retool
	retool do unparam $(PKGS)

errcheck: retool
	retool do errcheck $(PKGS)

lint: fmt vet megacheck misspell unconvert ineffassign unparam errcheck

# Building
build:
	go install github.com/jeromefroe/heimdallr/cmd/heimdallr

# Testing
test:
	go test -v ./...

test-full:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# CI
ci: lint test-full

# Release
release:
	@if test -z ${GITHUB_TOKEN}; then echo "GITHUB_TOKEN must be set to release a new version"; exit 1; fi
	git tag -a $(VERSION) -m "Releasing version $(VERSION)"
	git push origin $(VERSION)
	retool do goreleaser

# Docker
container:
	docker build -t $(IMAGE):$(VERSION) .

push: container
	docker push $(IMAGE):$(VERSION)
	@if git describe --tags --exact-match >/dev/null 2>&1; \
	then \
	    docker tag $(IMAGE):$(VERSION) $(IMAGE):latest; \
	    docker push $(IMAGE):latest; \
	fi
