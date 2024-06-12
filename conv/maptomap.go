// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conv

import (
	"reflect"

	"github.com/gocarp/codes"
	"github.com/gocarp/errors"
)

// MapToMap converts any map type variable `params` to another map type variable `pointer`
// using reflect.
// See doMapToMap.
func MapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return Scan(params, pointer, mapping...)
}

// doMapToMap converts any map type variable `params` to another map type variable `pointer`.
//
// The parameter `params` can be any type of map, like:
// map[string]string, map[string]struct, map[string]*struct, reflect.Value, etc.
//
// The parameter `pointer` should be type of *map, like:
// map[int]string, map[string]struct, map[string]*struct, reflect.Value, etc.
//
// The optional parameter `mapping` is used for struct attribute to map key mapping, which makes
// sense only if the items of original map `params` is type struct.
func doMapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var (
		paramsRv                  reflect.Value
		paramsKind                reflect.Kind
		keyToAttributeNameMapping map[string]string
	)
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	if v, ok := params.(reflect.Value); ok {
		paramsRv = v
	} else {
		paramsRv = reflect.ValueOf(params)
	}
	paramsKind = paramsRv.Kind()
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return doMapToMap(Map(params), pointer, mapping...)
	}
	// Empty params map, no need continue.
	if paramsRv.Len() == 0 {
		return nil
	}
	var pointerRv reflect.Value
	if v, ok := pointer.(reflect.Value); ok {
		pointerRv = v
	} else {
		pointerRv = reflect.ValueOf(pointer)
	}
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return errors.NewCodef(
			codes.CodeInvalidParameter,
			`destination pointer should be type of *map, but got: %s`,
			pointerKind,
		)
	}
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
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		pointerValueKind = pointerValueType.Kind()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	)
	// Retrieve the true element type of target map.
	if pointerValueKind == reflect.Ptr {
		pointerValueKind = pointerValueType.Elem().Kind()
	}
	for _, key := range paramsKeys {
		mapValue := reflect.New(pointerValueType).Elem()
		switch pointerValueKind {
		case reflect.Map, reflect.Struct:
			if err = doStruct(
				paramsRv.MapIndex(key).Interface(), mapValue, keyToAttributeNameMapping, "",
			); err != nil {
				return err
			}
		default:
			mapValue.Set(
				reflect.ValueOf(
					doConvert(doConvertInput{
						FromValue:  paramsRv.MapIndex(key).Interface(),
						ToTypeName: pointerValueType.String(),
						ReferValue: mapValue,
						Extra:      nil,
					}),
				),
			)
		}
		var mapKey = reflect.ValueOf(
			doConvert(doConvertInput{
				FromValue:  key.Interface(),
				ToTypeName: pointerKeyType.Name(),
				ReferValue: reflect.New(pointerKeyType).Elem().Interface(),
				Extra:      nil,
			}),
		)
		dataMap.SetMapIndex(mapKey, mapValue)
	}
	pointerRv.Set(dataMap)
	return nil
}
