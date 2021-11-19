# \AllPolicyTypeIdentifiersApi

All URIs are relative to *https://example.com/A1-P/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PolicytypesGet**](AllPolicyTypeIdentifiersApi.md#PolicytypesGet) | **Get** /policytypes | 



## PolicytypesGet

> []string PolicytypesGet(ctx).Execute()





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

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.AllPolicyTypeIdentifiersApi.PolicytypesGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AllPolicyTypeIdentifiersApi.PolicytypesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PolicytypesGet`: []string
    fmt.Fprintf(os.Stdout, "Response from `AllPolicyTypeIdentifiersApi.PolicytypesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiPolicytypesGetRequest struct via the builder pattern


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

