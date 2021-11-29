#!/bin/bash

## Verifying that the specification is correct first
#ToDo - enable verification through Docker container
#docker run -v api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml:/a1ap_enrichment_information.yaml --rm p1c2u/openapi-spec-validator /a1ap_enrichment_information.yaml
#docker run -v api/northbound/v301/policy_management/a1ap_policy_management.yaml:/a1ap_policy_management.yaml --rm p1c2u/openapi-spec-validator /a1ap_policy_management.yaml

openapi-spec-validator api/northbound/v301/policy_management/a1ap_policy_management.yaml
openapi-spec-validator api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml

## Building the code out of specification

# ToDo - bring the Docker build in through the Makefile
#docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
#    -i /local/api/northbound/v301/enrichment_information/a1ap_enrichment_information.json \
#    -g go \
#    -o /local/pkg/northbound/a1ap
#    -i https://raw.githubusercontent.com/openapitools/openapi-generator/master/modules/openapi-generator/src/test/resources/3_0/petstore.yaml \

# OpenAPI generator doesn't work for A1AP Policy Management even if validation is passes..
#openapi-generator generate -i api/northbound/v301/policy_management/a1ap_policy_management.yaml \
#                               -g go \
#                               -o pkg/northbound/a1ap

#openapi-generator generate -i api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml \
#                              -g go \
#                              -o pkg/northbound/a1ap

# Alternative way, which works so far..
# ToDo - create a Docker image to wrap up code generation with oapi-codegen
oapi-codegen api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml > pkg/northbound/a1ap/enrichment_information/a1ap_ei.go
oapi-codegen api/northbound/v301/policy_management/a1ap_policy_management.yaml > pkg/northbound/a1ap/policy_management/a1ap_pm.go
