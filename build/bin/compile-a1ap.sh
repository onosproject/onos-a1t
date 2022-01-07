#!/bin/bash

export ONOS_ROOT=$GOPATH/src/github.com/onosproject
export OAPI_SPEC_VALIDATOR_VERSION=0.3.1
export OAPI_CODEGEN_VERSION=1.6.6

## Validating specification first
docker run -v ${ONOS_ROOT}/onos-a1t/api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml:/openapi.yaml --rm p1c2u/openapi-spec-validator:${OAPI_SPEC_VALIDATOR_VERSION} /openapi.yaml
docker run -v ${ONOS_ROOT}/onos-a1t/api/northbound/v301/policy_management/a1ap_policy_management.yaml:/openapi.yaml --rm p1c2u/openapi-spec-validator:${OAPI_SPEC_VALIDATOR_VERSION} /openapi.yaml

# Old way, requires installing prerequisites with make openapi-spec-validator
#openapi-spec-validator api/northbound/v301/policy_management/a1ap_policy_management.yaml
#openapi-spec-validator api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml

## Building the code out of specification
#docker run -v ${ONOS_ROOT}/onos-a1t:${ONOS_ROOT}/onos-a1tapi/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml \
#  -w ${ONOS_ROOT}/onos-a1t:${ONOS_ROOT}/onos-a1t/pkg/northbound/a1ap/enrichment_information/a1ap_eii.go \
#	--rm nostrict/oapi-codegen:${OAPI_CODEGEN_VERSION}

# Old way, requires installing prerequisites with make oapi-codegen
oapi-codegen api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml > pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
oapi-codegen api/northbound/v301/policy_management/a1ap_policy_management.yaml > pkg/northbound/a1ap/policy_management/a1ap_pm.go
