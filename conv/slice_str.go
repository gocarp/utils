// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conv

import (
	"reflect"

	"github.com/gocarp/helpers/json"
	"github.com/gocarp/helpers/reflection"
)

// SliceStr is alias of Strings.
func SliceStr(any interface{}) []string {
	return Strings(any)
}

// Strings converts `any` to []string.
func Strings(any interface{}) []string {
	if any == nil {
		return nil
	}
	var (
		array []string = nil
	)
	switch value := any.(type) {
	case []int:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int8:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []int64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint8:
		if json.Valid(value) {
			_ = json.UnmarshalUseNumber(value, &array)
		}
		if array == nil {
			array = make([]string, len(value))
			for k, v := range value {
				array[k] = String(v)
			}
			return array
		}
	case string:
		byteValue := []byte(value)
		if json.Valid(byteValue) {
			_ = json.UnmarshalUseNumber(byteValue, &array)
		}
		if array == nil {
			if value == "" {
				return []string{}
			}
			// Prevent strings from being null
			// See Issue 3465 for details
			return []string{value}
		}
	case []uint16:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []uint64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []bool:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float32:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []float64:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []interface{}:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	case []string:
		array = value
	case [][]byte:
		array = make([]string, len(value))
		for k, v := range value {
			array[k] = String(v)
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(iStrings); ok {
		return v.Strings()
	}
	if v, ok := any.(iInterfaces); ok {
		return Strings(v.Interfaces())
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]string, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = String(originValueAndKind.OriginValue.Index(i).Interface())
		}
		return slice

	default:
		if originValueAndKind.OriginValue.IsZero() {
			return []string{}
		}
		return []string{String(any)}
	}
}
