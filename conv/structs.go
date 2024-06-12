// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conv

import (
	"reflect"

	"github.com/gocarp/codes"
	"github.com/gocarp/errors"
)

// Structs converts any slice to given struct slice.
// Also see Scan, Struct.
func Structs(params interface{}, pointer interface{}, paramKeyToAttrMap ...map[string]string) (err error) {
	return Scan(params, pointer, paramKeyToAttrMap...)
}

// SliceStruct is alias of Structs.
func SliceStruct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return Structs(params, pointer, mapping...)
}

// StructsTag acts as Structs but also with support for priority tag feature, which retrieves the
// specified tags for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple tags that can be joined with char ','.
func StructsTag(params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStructs(params, pointer, nil, priorityTag)
}

// doStructs converts any slice to given struct slice.
//
// It automatically checks and converts json string to []map if `params` is string/[]byte.
//
// The parameter `pointer` should be type of pointer to slice of struct.
// Note that if `pointer` is a pointer to another pointer of type of slice of struct,
// it will create the struct/pointer internally.
func doStructs(
	params interface{}, pointer interface{}, paramKeyToAttrMap map[string]string, priorityTag string,
) (err error) {
	defer func() {
		// Catch the panic, especially the reflection operation panics.
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && errors.HasStack(v) {
				err = v
			} else {
				err = errors.NewCodeSkipf(codes.CodeInternalPanic, 1, "%+v", exception)
			}
		}
	}()

	// Pointer type check.
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return errors.NewCodef(
				codes.CodeInvalidParameter,
				"pointer should be type of pointer, but got: %v", kind,
			)
		}
	}
	// Converting `params` to map slice.
	var (
		paramsList []interface{}
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	for paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	switch paramsKind {
	case reflect.Slice, reflect.Array:
		paramsList = make([]interface{}, paramsRv.Len())
		for i := 0; i < paramsRv.Len(); i++ {
			paramsList[i] = paramsRv.Index(i).Interface()
		}
	default:
		var paramsMaps = Maps(params)
		paramsList = make([]interface{}, len(paramsMaps))
		for i := 0; i < len(paramsMaps); i++ {
			paramsList[i] = paramsMaps[i]
		}
	}
	// If `params` is an empty slice, no conversion.
	if len(paramsList) == 0 {
		return nil
	}
	var (
		reflectElemArray = reflect.MakeSlice(pointerRv.Type().Elem(), len(paramsList), len(paramsList))
		itemType         = reflectElemArray.Index(0).Type()
		itemTypeKind     = itemType.Kind()
		pointerRvElem    = pointerRv.Elem()
		pointerRvLength  = pointerRvElem.Len()
	)
	if itemTypeKind == reflect.Ptr {
		// Pointer element.
		for i := 0; i < len(paramsList); i++ {
			var tempReflectValue reflect.Value
			if i < pointerRvLength {
				// Might be nil.
				tempReflectValue = pointerRvElem.Index(i).Elem()
			}
			if !tempReflectValue.IsValid() {
				tempReflectValue = reflect.New(itemType.Elem()).Elem()
			}
			if err = doStruct(paramsList[i], tempReflectValue, paramKeyToAttrMap, priorityTag); err != nil {
				return err
			}
			reflectElemArray.Index(i).Set(tempReflectValue.Addr())
		}
	} else {
		// Struct element.
		for i := 0; i < len(paramsList); i++ {
			var tempReflectValue reflect.Value
			if i < pointerRvLength {
				tempReflectValue = pointerRvElem.Index(i)
			} else {
				tempReflectValue = reflect.New(itemType).Elem()
			}
			if err = doStruct(paramsList[i], tempReflectValue, paramKeyToAttrMap, priorityTag); err != nil {
				return err
			}
			reflectElemArray.Index(i).Set(tempReflectValue)
		}
	}
	pointerRv.Elem().Set(reflectElemArray)
	return nil
}
