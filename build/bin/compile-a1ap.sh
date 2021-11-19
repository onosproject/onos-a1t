#!/bin/bash



#docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
#    -i /local/api/northbound/v301/enrichment_information/a1ap_enrichment_information.json \
#    -g go \
#    -o /local/pkg/northbound/a1ap
#    -i https://raw.githubusercontent.com/openapitools/openapi-generator/master/modules/openapi-generator/src/test/resources/3_0/petstore.yaml \

openapi-generator generate -i api/northbound/v301/policy_management/a1ap_policy_management.yaml \
                               -g go \
                               -o /local/pkg/northbound/a1ap

#openapi-generator generate -i api/northbound/v301/enrichment_information/a1ap_enrichment_information.yaml \
#                              -g go \
#                              -o /local/pkg/northbound/a1ap