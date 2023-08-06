/*
ChromaDB API

This is OpenAPI schema for ChromaDB API.

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// MetadatasInnerValue struct for MetadatasInnerValue
type MetadatasInnerValue struct {
	bool    *bool
	float32 *float32
	int32   *int32
	string  *string
}

// Unmarshal JSON data into any of the pointers in the struct
func (dst *MetadatasInnerValue) UnmarshalJSON(data []byte) error {
	var err error
	// try to unmarshal JSON data into bool
	err = json.Unmarshal(data, &dst.bool)
	if err == nil {
		jsonbool, _ := json.Marshal(dst.bool)
		if string(jsonbool) == "{}" { // empty struct
			dst.bool = nil
		} else {
			return nil // data stored in dst.bool, return on the first match
		}
	} else {
		dst.bool = nil
	}

	// try to unmarshal JSON data into float32
	err = json.Unmarshal(data, &dst.float32)
	if err == nil {
		jsonfloat32, _ := json.Marshal(dst.float32)
		if string(jsonfloat32) == "{}" { // empty struct
			dst.float32 = nil
		} else {
			return nil // data stored in dst.float32, return on the first match
		}
	} else {
		dst.float32 = nil
	}

	// try to unmarshal JSON data into int32
	err = json.Unmarshal(data, &dst.int32)
	if err == nil {
		jsonint32, _ := json.Marshal(dst.int32)
		if string(jsonint32) == "{}" { // empty struct
			dst.int32 = nil
		} else {
			return nil // data stored in dst.int32, return on the first match
		}
	} else {
		dst.int32 = nil
	}

	// try to unmarshal JSON data into string
	err = json.Unmarshal(data, &dst.string)
	if err == nil {
		jsonstring, _ := json.Marshal(dst.string)
		if string(jsonstring) == "{}" { // empty struct
			dst.string = nil
		} else {
			return nil // data stored in dst.string, return on the first match
		}
	} else {
		dst.string = nil
	}

	return fmt.Errorf("data failed to match schemas in anyOf(MetadatasInnerValue)")
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src *MetadatasInnerValue) MarshalJSON() ([]byte, error) {
	if src.bool != nil {
		return json.Marshal(&src.bool)
	}

	if src.float32 != nil {
		return json.Marshal(&src.float32)
	}

	if src.int32 != nil {
		return json.Marshal(&src.int32)
	}

	if src.string != nil {
		return json.Marshal(&src.string)
	}

	return nil, nil // no data in anyOf schemas
}

type NullableMetadatasInnerValue struct {
	value *MetadatasInnerValue
	isSet bool
}

func (v NullableMetadatasInnerValue) Get() *MetadatasInnerValue {
	return v.value
}

func (v *NullableMetadatasInnerValue) Set(val *MetadatasInnerValue) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadatasInnerValue) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadatasInnerValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadatasInnerValue(val *MetadatasInnerValue) *NullableMetadatasInnerValue {
	return &NullableMetadatasInnerValue{value: val, isSet: true}
}

func (v NullableMetadatasInnerValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadatasInnerValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
