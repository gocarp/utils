// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conv

import (
	"time"

	"github.com/gocarp/go/times"
	"github.com/gocarp/helpers/utils"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := GTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := times.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// GTime converts `any` to *times.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using times.NewFromTimeStamp if `any` is numeric,
// or using times.StrToTime if `any` is string.
func GTime(any interface{}, format ...string) *times.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(iGTime); ok {
		return v.GTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*times.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return times.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return times.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return times.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := times.StrToTimeFormat(s, item)
			if t != nil && err == nil {
				return t
			}
		}
		return nil
	}
	if utils.IsNumeric(s) {
		return times.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := times.StrToTime(s)
		return t
	}
}
