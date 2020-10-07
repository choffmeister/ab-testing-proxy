.PHONY: run test test-watch build push backends

# The binary to build (just the basename).
BIN := ab-testing-proxy

# Where to push the docker image.
REGISTRY ?= choffmeister

IMAGE := $(REGISTRY)/$(BIN)-amd64

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

run: backends
	cd src && go run . --config ../example.config.yaml

test: backends
	cd src && go test -v

test-watch: backends
	cd src && watch -n1 go test -v

build:
	docker build -t $(IMAGE):$(VERSION) .

push: build
	docker push $(IMAGE):$(VERSION)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):latest

backends:
	docker-compose up -d
	sleep 1
