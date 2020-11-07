package common

import (
	"encoding/json"
	"fmt"
	"reflect"

	. "github.com/fhivemind/go-hastily/pkg/global"
)

// Generic is a generic object definition which maps
// struct public fields to its keys and values.
type Generic struct {
	Keys   []string
	Values []string
}

// ObjectToGeneric creates a generic object from an interface.
func ObjectToGeneric(object interface{}) *Generic {
	var result Generic
	result.Keys, result.Values = StructKeysAndValues(object)
	return &result
}

// ObjectToMap converts interface to map of string keys.
// Note: The object should have JSON conversion defined.
func ObjectToMap(object interface{}) map[string]interface{} {
	var objMap map[string]interface{}
	toJson, _ := json.Marshal(object)
	json.Unmarshal(toJson, &objMap)
	return objMap
}

// StructKeysAndValues extracts struct field names and values into arrays of strings.
func StructKeysAndValues(data interface{}) ([]string, []string) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var keys []string
	var values []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanSet() {
			keys = append(keys, v.Type().Field(i).Name)
			values = append(values, fmt.Sprintf("%+v", v.Field(i).Interface()))
		}
	}

	return keys, values
}

// StructNonNullKeysAndValues extracts struct field names and values into arrays of strings.
func StructNonNullKeysAndValues(data interface{}) ([]string, []string) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var keys []string
	var values []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanSet() {
			key := v.Type().Field(i).Name
			value := v.Field(i).Interface()
			if !IsZero(value) {
				keys = append(keys, key)
				values = append(values, fmt.Sprintf("%+v", value))
			}
		}
	}

	return keys, values
}
