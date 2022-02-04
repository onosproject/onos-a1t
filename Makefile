# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

.PHONY: build

OAPI_CODEGEN_VERSION := v1.9.0
OAPI_SPEC_VALIDATOR_VERSION := 0.3.1

ONOS_A1T_VERSION := latest
ONOS_BUILD_VERSION := v0.6.6
ONOS_PROTOC_VERSION := v0.6.6
BUF_VERSION := 0.27.1

build: # @HELP build the Go binaries and run all validations (default)
build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/onos-a1t ./cmd/onos-a1t

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

#ToDo - run it through Docker container in the future
build_api:
	build/bin/compile-a1ap.sh

# Requires providing a filename
oapi-codegen:
	oapi-codegen || ( cd .. && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION})

# Requires providing a filename
openapi-spec-validator:
	openapi-spec-validator || ( cd .. && pip3 install openapi-spec-validator==${OAPI_SPEC_VALIDATOR_VERSION})

license_check_a1t:  # @HELP examine and ensure license headers exist
	./build/build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR} --skipped-dir pkg/northbound/a1ap --skipped-dir api --skipped-dir build --boilerplate SPDX-Apache-2.0

buflint: #@HELP run the "buf check lint" command on the proto files in 'api'
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-a1t \
		-w /go/src/github.com/onosproject/onos-a1t/api \
		bufbuild/buf:${BUF_VERSION} check lint

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps linters license_check_a1t
	go test -race github.com/onosproject/onos-a1t/pkg/...
	go test -race github.com/onosproject/onos-a1t/cmd/...

jenkins-test: # @HELP run the unit tests and source code validation producing a junit style report for Jenkins
jenkins-test: build deps license_check_a1t linters
	TEST_PACKAGES=github.com/onosproject/onos-a1t/... ./build/build-tools/build/jenkins/make-unit

onos-a1t-docker: # @HELP build onos-a1t Docker image
onos-a1t-docker:
	@go mod vendor
	docker build . -f build/onos-a1t/Dockerfile \
		-t onosproject/onos-a1t:${ONOS_A1T_VERSION}
	@rm -rf vendor

images: # @HELP build all Docker images
images: build onos-a1t-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-a1t:${ONOS_A1T_VERSION}

all: build images

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION} onosproject/onos-a1t

jenkins-publish: jenkins-tools # @HELP Jenkins calls this to publish artifacts
	./build/bin/push-images
	./build/build-tools/release-merge-commit


protos: # @HELP build a1t golang protobuffer (TEMP) # TODO move .proto to onos-api
protos:
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/gogoproto
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/protoc-gen-gogofaster
	go_import_paths="Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types"
	proto_path="./api/southbound/:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"
	protoc --proto_path=$proto_path --gogofaster_out=$go_import_paths,import_path=onos/a1t,plugins=grpc:./pkg/southbound/a1t/ a1t.proto