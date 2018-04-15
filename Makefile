# Go build flags
GOOS := linux
GOARCH := amd64

REGISTRY := quay.io/munnerz
APP_NAME := apiextensions-ca-helper
IMAGE_TAGS := canary

BUILD_TAG := build

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-tags netgo \
		-o ${APP_NAME}-$@_$(GOOS)_$(GOARCH) \
		.

docker_build:
	docker build \
		-t $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) \
		-f Dockerfile \
		.

docker_push:
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker tag $(REGISTRY)/$(APP_NAME):$(BUILD_TAG) $(REGISTRY)/$(APP_NAME):$${tag} ; \
		docker push $(REGISTRY)/$(APP_NAME):$${tag}; \
	done
