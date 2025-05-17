# Set default shell to bash
SHELL := /bin/bash

# Determine OS and architecture
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
VERSION := 0.1.0

# Provider info
PROVIDER_NAME := terraform-provider-litellm
PROVIDER_PATH := registry.terraform.io/jimmyflatting/litellm
BUILD_DIR := bin

.PHONY: build
build:
	mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${PROVIDER_NAME}

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${PROVIDER_PATH}/${VERSION}/${OS}_${ARCH}
	cp ${BUILD_DIR}/${PROVIDER_NAME} ~/.terraform.d/plugins/${PROVIDER_PATH}/${VERSION}/${OS}_${ARCH}/

.PHONY: test
test:
	go test ./... -v

.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: all
all: fmt vet test build