# \AllPolicyIdentifiersApi

All URIs are relative to *https://example.com/A1-P/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PolicytypesPolicyTypeIdPoliciesGet**](AllPolicyIdentifiersApi.md#PolicytypesPolicyTypeIdPoliciesGet) | **Get** /policytypes/{policyTypeId}/policies | 



## PolicytypesPolicyTypeIdPoliciesGet

> []string PolicytypesPolicyTypeIdPoliciesGet(ctx, policyTypeId).Execute()





### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    policyTypeId := "policyTypeId_example" // string | 

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.AllPolicyIdentifiersApi.PolicytypesPolicyTypeIdPoliciesGet(context.Background(), policyTypeId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AllPolicyIdentifiersApi.PolicytypesPolicyTypeIdPoliciesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PolicytypesPolicyTypeIdPoliciesGet`: []string
    fmt.Fprintf(os.Stdout, "Response from `AllPolicyIdentifiersApi.PolicytypesPolicyTypeIdPoliciesGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**policyTypeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPolicytypesPolicyTypeIdPoliciesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**[]string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

