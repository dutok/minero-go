package nbt

import (
	"fmt"
	"io"

	"github.com/minero/minero/types"
)

func NewByteArray(s []int8) *ByteArray {
	b := &ByteArray{
		Value: make([]types.Int8, len(s)),
	}
	for index, elem := range s {
		b.Value[index] = types.Int8(elem)
	}
	return b
}

func NewIntArray(s []int32) *IntArray {
	b := &IntArray{
		Value: make([]types.Int32, len(s)),
	}
	for index, elem := range s {
		b.Value[index] = types.Int32(elem)
	}
	return b
}

// ByteArray holds a length-prefixed array of signed bytes. The prefix is a
// signed integer (4 bytes).
// TagType: 7, Size: 4 + elem * 1 bytes
type ByteArray struct {
	Value []types.Int8
}

func (arr ByteArray) Type() TagType          { return TagByteArray }
func (arr ByteArray) Size() int64            { return int64(4 + len(arr.Value)) }
func (arr ByteArray) Lookup(path string) Tag { return nil }

func (arr ByteArray) String() string {
	return fmt.Sprintf("Array(%d x Byte)", len(arr.Value))
}

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (arr *ByteArray) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read length-prefix
	var length Int32
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length bytes
	arr.Value = make([]types.Int8, length.Int32)
	for index, elem := range arr.Value {
		nn, err = elem.ReadFrom(r)
		if err != nil {
			return
		}
		arr.Value[index] = elem
		n += nn
	}

	return
}

// WriteTo satifies io.WriterTo interface. TypeId is not encoded.
func (arr *ByteArray) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Write length-prefix
	var length = types.Int32(len(arr.Value))
	nn, err = length.WriteTo(w)
	if err != nil {
		return
	}
	n += nn

	// Then write byte array
	for _, elem := range arr.Value {
		nn, err = elem.WriteTo(w)
		if err != nil {
			return
		}
		n += nn
	}

	return
}

// IntArray holds a length-prefixed array of signed integers. The prefix is a
// signed integer (4 bytes) and indicates the number of 4 byte integers.
// TagType: 11, Size: 4 + 4 * elem
type IntArray struct {
	Value []types.Int32
}

func (arr IntArray) Type() TagType          { return TagIntArray }
func (arr IntArray) Size() int64            { return int64(4 + len(arr.Value)) }
func (arr IntArray) Lookup(path string) Tag { return nil }
func (arr IntArray) String() string {
	return fmt.Sprintf("Array(%d x Int)", len(arr.Value))
}

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (arr *IntArray) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read length-prefix
	var length Int32
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length bytes
	arr.Value = make([]types.Int32, length.Int32)
	for index, elem := range arr.Value {
		nn, err = elem.ReadFrom(r)
		if err != nil {
			return
		}
		arr.Value[index] = elem
		n += nn
	}

	return
}

// WriteTo satifies io.WriterTo interface. TypeId is not encoded.
func (arr *IntArray) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Write length-prefix
	var length = types.Int32(len(arr.Value))
	nn, err = length.WriteTo(w)
	if err != nil {
		return
	}
	n += nn

	// Then write int array
	for _, tag := range arr.Value {
		nn, err = tag.WriteTo(w)
		if err != nil {
			return
		}
		n += nn
	}
	return
}
