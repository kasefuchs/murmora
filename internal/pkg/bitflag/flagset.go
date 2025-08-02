// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package bitflag

import (
	"database/sql/driver"
	"fmt"
)

type FlagSet[T BitwiseEnum] struct {
	raw int64
}

func NewFlagSet[T BitwiseEnum]() *FlagSet[T] {
	return &FlagSet[T]{}
}

func (f *FlagSet[T]) Raw() int64 {
	return f.raw
}

func (f *FlagSet[T]) SetRaw(raw int64) {
	f.raw = raw
}

func (f *FlagSet[T]) Add(flag T) {
	f.raw |= BitFromEnum(flag)
}

func (f *FlagSet[T]) Remove(flag T) {
	f.raw &^= BitFromEnum(flag)
}

func (f *FlagSet[T]) Has(flag T) bool {
	return f.raw&BitFromEnum(flag) != 0
}

func (f *FlagSet[T]) Clear() {
	f.raw = 0
}

func (f *FlagSet[T]) Value() (driver.Value, error) {
	return f.raw, nil
}

func (f *FlagSet[T]) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		f.raw = v
	case int32:
		f.raw = int64(v)
	case int:
		f.raw = int64(v)
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
	return nil
}
