// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"reflect"

	"github.com/gocarp/helpers/empty"
)

// IsEmpty checks given `value` empty or not.
// It returns false if `value` is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}

// IsTypeOf checks and returns whether the type of `value` and `valueInExpectType` equal.
func IsTypeOf(value, valueInExpectType interface{}) bool {
	return reflect.TypeOf(value) == reflect.TypeOf(valueInExpectType)
}
