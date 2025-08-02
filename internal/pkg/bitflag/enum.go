// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package bitflag

type BitwiseEnum interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

func BitFromEnum[T BitwiseEnum](e T) int64 {
	val := int64(e)
	if val <= 0 {
		return 0
	}

	return 1 << (val - 1)
}
