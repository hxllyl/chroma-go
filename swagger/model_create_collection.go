/*
ChromaDB API

This is OpenAPI schema for ChromaDB API.

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the CreateCollection type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateCollection{}

// CreateCollection struct for CreateCollection
type CreateCollection struct {
	Name        string                 `json:"name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	GetOrCreate *bool                  `json:"get_or_create,omitempty"`
}

// NewCreateCollection instantiates a new CreateCollection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateCollection(name string) *CreateCollection {
	this := CreateCollection{}
	this.Name = name
	var getOrCreate bool = false
	this.GetOrCreate = &getOrCreate
	return &this
}

// NewCreateCollectionWithDefaults instantiates a new CreateCollection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateCollectionWithDefaults() *CreateCollection {
	this := CreateCollection{}
	var getOrCreate bool = false
	this.GetOrCreate = &getOrCreate
	return &this
}

// GetName returns the Name field value
func (o *CreateCollection) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CreateCollection) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CreateCollection) SetName(v string) {
	o.Name = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *CreateCollection) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateCollection) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *CreateCollection) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *CreateCollection) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetGetOrCreate returns the GetOrCreate field value if set, zero value otherwise.
func (o *CreateCollection) GetGetOrCreate() bool {
	if o == nil || IsNil(o.GetOrCreate) {
		var ret bool
		return ret
	}
	return *o.GetOrCreate
}

// GetGetOrCreateOk returns a tuple with the GetOrCreate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateCollection) GetGetOrCreateOk() (*bool, bool) {
	if o == nil || IsNil(o.GetOrCreate) {
		return nil, false
	}
	return o.GetOrCreate, true
}

// HasGetOrCreate returns a boolean if a field has been set.
func (o *CreateCollection) HasGetOrCreate() bool {
	if o != nil && !IsNil(o.GetOrCreate) {
		return true
	}

	return false
}

// SetGetOrCreate gets a reference to the given bool and assigns it to the GetOrCreate field.
func (o *CreateCollection) SetGetOrCreate(v bool) {
	o.GetOrCreate = &v
}

func (o CreateCollection) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateCollection) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsNil(o.GetOrCreate) {
		toSerialize["get_or_create"] = o.GetOrCreate
	}
	return toSerialize, nil
}

type NullableCreateCollection struct {
	value *CreateCollection
	isSet bool
}

func (v NullableCreateCollection) Get() *CreateCollection {
	return v.value
}

func (v *NullableCreateCollection) Set(val *CreateCollection) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateCollection) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateCollection) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateCollection(val *CreateCollection) *NullableCreateCollection {
	return &NullableCreateCollection{value: val, isSet: true}
}

func (v NullableCreateCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateCollection) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
