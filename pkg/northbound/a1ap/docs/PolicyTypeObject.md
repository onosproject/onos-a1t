# PolicyTypeObject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PolicySchema** | **map[string]interface{}** | A JSON schema following http://json-schema.org/draft-07/schema | 
**StatusSchema** | Pointer to **map[string]interface{}** | A JSON schema following http://json-schema.org/draft-07/schema | [optional] 

## Methods

### NewPolicyTypeObject

`func NewPolicyTypeObject(policySchema map[string]interface{}, ) *PolicyTypeObject`

NewPolicyTypeObject instantiates a new PolicyTypeObject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPolicyTypeObjectWithDefaults

`func NewPolicyTypeObjectWithDefaults() *PolicyTypeObject`

NewPolicyTypeObjectWithDefaults instantiates a new PolicyTypeObject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPolicySchema

`func (o *PolicyTypeObject) GetPolicySchema() map[string]interface{}`

GetPolicySchema returns the PolicySchema field if non-nil, zero value otherwise.

### GetPolicySchemaOk

`func (o *PolicyTypeObject) GetPolicySchemaOk() (*map[string]interface{}, bool)`

GetPolicySchemaOk returns a tuple with the PolicySchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicySchema

`func (o *PolicyTypeObject) SetPolicySchema(v map[string]interface{})`

SetPolicySchema sets PolicySchema field to given value.


### GetStatusSchema

`func (o *PolicyTypeObject) GetStatusSchema() map[string]interface{}`

GetStatusSchema returns the StatusSchema field if non-nil, zero value otherwise.

### GetStatusSchemaOk

`func (o *PolicyTypeObject) GetStatusSchemaOk() (*map[string]interface{}, bool)`

GetStatusSchemaOk returns a tuple with the StatusSchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusSchema

`func (o *PolicyTypeObject) SetStatusSchema(v map[string]interface{})`

SetStatusSchema sets StatusSchema field to given value.

### HasStatusSchema

`func (o *PolicyTypeObject) HasStatusSchema() bool`

HasStatusSchema returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


