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

#ToDo - run it through Docker container in the future
build_api:
	build/bin/compile-a1ap.sh

# Requires providing a filename
oapi-codegen:
	oapi-codegen || ( cd .. && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION})

# Requires providing a filename
openapi-spec-validator:
	openapi-spec-validator || ( cd .. && pip3 install openapi-spec-validator==${OAPI_SPEC_VALIDATOR_VERSION})

build-tools: # @HELP install the ONOS build tools if needed
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi

jenkins-tools: # @HELP installs tooling needed for Jenkins
	cd .. && go get -u github.com/jstemmer/go-junit-report && go get github.com/t-yuki/gocover-cobertura

deps: # @HELP ensure that the required dependencies are in place
	GOPRIVATE="github.com/onosproject/*" go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

license_check: build-tools # @HELP examine and ensure license headers exist
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR} --boilerplate LicenseRef-ONF-Member-Only-1.0

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 10m

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.42.0

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

buflint: #@HELP run the "buf check lint" command on the proto files in 'api'
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-a1t \
		-w /go/src/github.com/onosproject/onos-a1t/api \
		bufbuild/buf:${BUF_VERSION} check lint

test: # @HELP run the unit tests and source code validation producing a golang style report
test: build deps linters license_check
	go test -race github.com/onosproject/onos-a1t/...

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

protos: # @HELP build a1t golang protobuffer (TEMP) # TODO move .proto to onos-api
protos:
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/gogoproto
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/protoc-gen-gogofaster
	go_import_paths="Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types"
	proto_path="./api/southbound/:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"
	protoc --proto_path=$proto_path --gogofaster_out=$go_import_paths,import_path=onos/a1t,plugins=grpc:./pkg/southbound/a1t/ a1t.proto