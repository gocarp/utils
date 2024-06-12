// Copyright (c) 2022-2024 The Focela Authors, All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import "github.com/gocarp/helpers/deepcopy"

// Copy returns a deep copy of v.
//
// Copy is unable to copy unexported fields in a struct (lowercase field names).
// Unexported fields can't be reflected by the Go runtime, and therefore
// they can't perform any data copies.
func Copy(src interface{}) (dst interface{}) {
	return deepcopy.Copy(src)
}
