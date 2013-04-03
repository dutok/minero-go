package nbt

import (
	"fmt"
	"io"

	"github.com/toqueteos/minero/types"
)

// ByteArray holds a length-prefixed array of signed bytes. The prefix is a
// signed integer (4 bytes).
// TagType: 7, Size: 4 + elem * 1 bytes
type ByteArray struct {
	Value []types.Byte
}

func (b ByteArray) Type() TagType          { return TagByteArray }
func (b ByteArray) Size() int64            { return int64(4 + len(b.Value)) }
func (b ByteArray) Lookup(path string) Tag { return nil }

func (b ByteArray) String() string {
	s := "NBT_ByteArray('%s', size: %d) [% x, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(b.Value), b.Value[:ArrayNum], len(b.Value)-ArrayNum)
}

func (b *ByteArray) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read length-prefix
	var length Int
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length bytes
	arr := make([]types.Byte, length.Int)
	for _, elem := range arr {
		nn, err = elem.ReadFrom(r)
		if err != nil {
			return
		}
		n += nn
	}

	b.Value = arr
	return
}

func (b *ByteArray) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Write length-prefix
	var length = types.Int(len(b.Value))
	nn, err = length.WriteTo(w)
	if err != nil {
		return
	}
	n += nn

	// Then write byte array
	for _, elem := range b.Value {
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
	Value []types.Int
}

func (i IntArray) Type() TagType          { return TagIntArray }
func (i IntArray) Size() int64            { return int64(4 + len(i.Value)) }
func (i IntArray) Lookup(path string) Tag { return nil }
func (i IntArray) String() string {
	var values []int32
	for _, elem := range i.Value[:ArrayNum] {
		values = append(values, int32(elem))
	}

	s := "NBT_IntArray('%s', size: %d) [% d, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(i.Value), values, len(i.Value)-ArrayNum)
}

func (i *IntArray) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read length-prefix
	var length Int
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length bytes
	arr := make([]types.Int, length.Int)
	for _, elem := range arr {
		nn, err = elem.ReadFrom(r)
		if err != nil {
			return
		}
		n += nn
	}

	i.Value = arr
	return
}

func (i *IntArray) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Write length-prefix
	var length = types.Int(len(i.Value))
	nn, err = length.WriteTo(w)
	if err != nil {
		return
	}
	n += nn

	// Then write int array
	for _, tag := range i.Value {
		nn, err = tag.WriteTo(w)
		if err != nil {
			return
		}
		n += nn
	}
	return
}
