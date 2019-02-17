# Go build flags
GOOS := linux
ALL_ARCHES:=amd64 arm arm64

REGISTRY := quay.io/munnerz
APP_NAME := apiextensions-ca-helper
DOCKER_IMAGE := $(REGISTRY)/$(APP_NAME)
IMAGE_TAGS := canary

BUILD_TAG := build

ifeq ($(ARCH),amd64)
	BASEIMAGE?=alpine:3.9
endif
ifeq ($(ARCH),arm)
	BASEIMAGE?=arm32v6/alpine:3.9
endif
ifeq ($(ARCH),arm64)
	BASEIMAGE?=arm64v8/alpine:3.9
endif

all: clean build

build: $(addprefix build/apiextensions-ca-helper-,$(ALL_ARCHES))

.PHONY: clean
clean:
	go clean
	rm -rf ./build/

build/apiextensions-ca-helper-%: *.go
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$* go build -tags netgo -o $@ $^
	$(MAKE) ARCH=$* build-and-push

.PHONY: build-and-push
build-and-push:
	cat Dockerfile.build | sed "s|BASEIMAGE|$(BASEIMAGE)|g" | sed "s|ARCH|$(ARCH)|g" \
		> build/Dockerfile-$(ARCH)
	docker build -t $(DOCKER_IMAGE)-$(ARCH):$(BUILD_TAG) -f build/Dockerfile-$(ARCH) ./build
	@for tag in $(IMAGE_TAGS); do \
		docker tag $(DOCKER_IMAGE)-$(ARCH):$(BUILD_TAG) $(DOCKER_IMAGE)-$(ARCH):$${tag} ; \
		docker push $(DOCKER_IMAGE)-$(ARCH):$${tag}; \
	done

