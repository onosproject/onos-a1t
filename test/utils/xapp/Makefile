# SPDX-FileCopyrightText: 2022 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

.PHONY: build

ONOS_A1TXAPP_VERSION := latest

build: # @HELP build the Go binaries and run all validations (default)
build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/onos-a1txapp ./cmd/onos-a1txapp

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

onos-a1txapp-docker: # @HELP build onos-a1t Docker image
onos-a1txapp-docker:
	@go mod vendor
	docker build . -f build/onos-a1txapp/Dockerfile \
		-t onosproject/onos-a1txapp:${ONOS_A1TXAPP_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-a1txapp-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-a1txapp:${ONOS_A1TXAPP_VERSION}

all: build
