// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/kasefuchs/murmora/internal/pkg/bitflag"

func BitFieldFromFlagSet[T bitflag.BitwiseEnum](flagSet *bitflag.FlagSet[T]) *BitField {
	return &BitField{Value: flagSet.Raw()}
}

type TypedBitField[T bitflag.BitwiseEnum] struct {
	flagSet *bitflag.FlagSet[T]
}

func NewTypedBitField[T bitflag.BitwiseEnum](bitfield *BitField) *TypedBitField[T] {
	flagSet := bitflag.NewFlagSet[T]()

	if bitfield != nil {
		flagSet.SetRaw(bitfield.Value)
	}

	return &TypedBitField[T]{flagSet: flagSet}
}

func (t *TypedBitField[T]) ToProto() *BitField {
	return &BitField{
		Value: t.flagSet.Raw(),
	}
}

func (t *TypedBitField[T]) Raw() int64 {
	return t.flagSet.Raw()
}

func (t *TypedBitField[T]) SetRaw(raw int64) {
	t.flagSet.SetRaw(raw)
}

func (t *TypedBitField[T]) Add(flag T) {
	t.flagSet.Add(flag)
}

func (t *TypedBitField[T]) Remove(flag T) {
	t.flagSet.Remove(flag)
}

func (t *TypedBitField[T]) Has(flag T) bool {
	return t.flagSet.Has(flag)
}

func (t *TypedBitField[T]) Clear() {
	t.flagSet.Clear()
}
