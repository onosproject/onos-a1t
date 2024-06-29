# SPDX-License-Identifier: Apache-2.0
# Copyright 2019 Open Networking Foundation
# Copyright 2024 Intel Corporation

.PHONY: build

export CGO_ENABLED=1
export GO111MODULE=on

OAPI_CODEGEN_VERSION := v1.9.0
OAPI_SPEC_VALIDATOR_VERSION := 0.3.1

ONOS_A1T_VERSION ?= latest
ONOS_BUILD_VERSION := v0.6.6
ONOS_PROTOC_VERSION := v0.6.6
BUF_VERSION := 0.27.1

all: build docker-build

build: # @HELP build the Go binaries and run all validations (default)
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/onos-a1t ./cmd/onos-a1t

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build lint license
	go test -race github.com/onosproject/onos-a1t/pkg/...
	go test -race github.com/onosproject/onos-a1t/cmd/...

docker-build-onos-a1t: # @HELP build onos-a1t Docker image
	@go mod vendor
	docker build . -f build/onos-a1t/Dockerfile \
		-t onosproject/onos-a1t:${ONOS_A1T_VERSION}
	@rm -rf vendor

docker-build: # @HELP build all Docker images
docker-build: build docker-build-onos-a1t

docker-push-onos-a1t: # @HELP push onos-a1t Docker image
	docker push onosproject/onos-a1t:${ONOS_A1T_VERSION}

docker-push: # @HELP push docker images
docker-push: docker-push-onos-a1t

lint: # @HELP examines Go source code and reports coding problems
	golangci-lint --version | grep $(GOLANG_CI_VERSION) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin $(GOLANG_CI_VERSION)
	golangci-lint run --timeout 15m

license: # @HELP run license checks
	rm -rf venv
	python3 -m venv venv
	. ./venv/bin/activate;\
	python3 -m pip install --upgrade pip;\
	python3 -m pip install reuse;\
	reuse lint

check-version: # @HELP check version is duplicated
	./build/bin/version_check.sh all


#ToDo - run it through Docker container in the future
build-api:
	build/bin/compile-a1ap.sh

# Requires providing a filename
oapi-codegen:
	oapi-codegen || ( cd .. && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION})

# Requires providing a filename
openapi-spec-validator:
	openapi-spec-validator || ( cd .. && pip3 install openapi-spec-validator==${OAPI_SPEC_VALIDATOR_VERSION})

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-a1t/onos-a1t ./cmd/onos/onos venv
	go clean github.com/onosproject/onos-a1t/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '