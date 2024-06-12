// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"reflect"

	"github.com/gocarp/codes"
	"github.com/gocarp/errors"
	"github.com/gocarp/go/structs"
	"github.com/gocarp/utils/conv"
)

// StructToSlice converts struct to slice of which all keys and values are its items.
// Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func StructToSlice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		array := make([]interface{}, 0)
		// Note that, it uses the conv tag name instead of the attribute name if
		// the conv tag is fined in the struct attributes.
		for k, v := range conv.Map(reflectValue) {
			array = append(array, k)
			array = append(array, v)
		}
		return array
	}
	return nil
}

// FillStructWithDefault fills  attributes of pointed struct with tag value from `default/d` tag .
// The parameter `structPtr` should be either type of *struct/[]*struct.
func FillStructWithDefault(structPtr interface{}) error {
	var (
		reflectValue reflect.Value
	)
	if rv, ok := structPtr.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(structPtr)
	}
	switch reflectValue.Kind() {
	case reflect.Ptr:
		// Nothing to do.
	case reflect.Array, reflect.Slice:
		if reflectValue.Elem().Kind() != reflect.Ptr {
			return errors.NewCodef(
				codes.CodeInvalidParameter,
				`invalid parameter "%s", the element of slice should be type of pointer of struct, but given "%s"`,
				reflectValue.Type().String(), reflectValue.Elem().Type().String(),
			)
		}
	default:
		return errors.NewCodef(
			codes.CodeInvalidParameter,
			`invalid parameter "%s", should be type of pointer of struct`,
			reflectValue.Type().String(),
		)
	}
	if reflectValue.IsNil() {
		return errors.NewCode(
			codes.CodeInvalidParameter,
			`the pointed struct object should not be nil`,
		)
	}
	if !reflectValue.Elem().IsValid() {
		return errors.NewCode(
			codes.CodeInvalidParameter,
			`the pointed struct object should be valid`,
		)
	}
	fields, err := structs.Fields(structs.FieldsInput{
		Pointer:         reflectValue,
		RecursiveOption: structs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return err
	}
	for _, field := range fields {
		if field.OriginalKind() == reflect.Struct {
			err := FillStructWithDefault(field.OriginalValue().Addr())
			if err != nil {
				return err
			}
			continue
		}

		if defaultValue := field.TagDefault(); defaultValue != "" {
			if field.IsEmpty() {
				field.Value.Set(reflect.ValueOf(
					conv.ConvertWithRefer(defaultValue, field.Value),
				))
			}
		}
	}

	return nil
}
