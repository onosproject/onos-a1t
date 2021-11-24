OAPI_CODEGEN_VERSION := v1.9.0
OAPI_SPEC_VALIDATOR_VERSION := 0.3.1

.PHONY: build

build: # @HELP build the Go binaries and run all validations (default)
#build:
	#go build -o build/_output/onos-a1t ./cmd/onos-a1t

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

deps: # @HELP ensure that the required dependencies are in place
#	go build -v ./...
#	bash -c "diff -u <(echo -n) <(git diff go.mod)"
#	bash -c "diff -u <(echo -n) <(git diff go.sum)"

license_check: build-tools # @HELP examine and ensure license headers exist
	@#if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	#./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR} --boilerplate LicenseRef-ONF-Member-Only-1.0 --skipped-dir=python

linters: golang-ci # @HELP examines Go source code and reports coding problems
	golangci-lint run --timeout 10m

golang-ci: # @HELP install golang-ci if not present
	golangci-lint --version || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b `go env GOPATH`/bin v1.42.0

test: # @HELP run the unit tests and source code validation producing a golang style report
#test: build deps linters license_check
#	go test -race github.com/onosproject/onos-a1t/...
