.PHONY: build test release clean

GITHUB_USER := atsaki

VERSION ?= dev
TARGET_OS ?= darwin linux windows
TARGET_ARCH ?= amd64

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_PATH))))
PROJECT := github.com/$(GITHUB_USER)/$(CURRENT_DIR)

EXECUTABLE_NAME := $(CURRENT_DIR)
DOCKER_IMAGE_NAME := $(CURRENT_DIR)-build
DOCKER_CONTAINER_NAME := $(DOCKER_IMAGE_NAME)-build-container

default: build

clean:
	rm -rf dist

ifeq ($(USE_CONTAINER), true)
build test release:
	docker build -t $(DOCKER_IMAGE_NAME) .

	docker run \
		--rm \
		--name $(DOCKER_CONTAINER_NAME) \
		-v $(shell pwd):/go/src/$(PROJECT) \
		-w /go/src/$(PROJECT) \
		-e TARGET_OS \
		-e TARGET_ARCH \
		-e GITHUB_TOKEN \
		$(DOCKER_IMAGE_NAME) \
		make $@
else
build:
	glide --debug up

	GOGC=off gox -os "$(TARGET_OS)" -arch "$(TARGET_ARCH)" \
	-output "dist/$(EXECUTABLE_NAME)_{{.OS}}_{{.Arch}}" .

test:
	exit 0

release: clean build
	ghr --repository $(CURRENT_DIR) \
		--username $(GITHUB_USER) \
		--prerelease \
		--replace \
		$(VERSION) dist/
endif

