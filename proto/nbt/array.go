package nbt

import (
	"fmt"
	"io"

	"github.com/toqueteos/minero/types"
)

// ByteArray holds a length-prefixed array of signed bytes. The prefix is a
// signed integer (4 bytes).
// TagType: 7, Size: 4 + elem * 1 bytes
type ByteArray []types.Byte

func (*ByteArray) Type() TagType             { return TagByteArray }
func (b ByteArray) Size() int64              { return int64(4 + len(b)) }
func (*ByteArray) Lookup(path string) Tagger { return nil }

func (b ByteArray) String() string {
	s := "NBT_ByteArray(size: %d) [% x, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(b), b[:StringNum], len(b)-StringNum)
}

func (b *ByteArray) ReadFrom(reader io.Reader) (n int64, err error) {
	var length Int

	// Read length-prefix
	// _, err = length.ReadFrom(reader)
	if err != nil {
		return
	}

	// Read length bytes
	arr := make([]types.Byte, length.Int)
	for _, elem := range arr {
		_, err = elem.ReadFrom(reader)
		if err != nil {
			return
		}
	}
	// _, err = io.Copy(b, reader)

	*b = arr
	return
}

func (b *ByteArray) WriteTo(writer io.Writer) (n int64, err error) {
	var length = types.Int(len(*b))

	// Write length-prefix
	if _, err = length.WriteTo(writer); err != nil {
		return
	}

	// Then write byte array
	for _, elem := range *b {
		_, err = elem.WriteTo(writer)
		if err != nil {
			return
		}
	}
	// err = io.Copy(writer, b)

	return
}

// IntArray holds a length-prefixed array of signed integers. The prefix is a
// signed integer (4 bytes) and indicates the number of 4 byte integers.
// TagType: 11, Size: 4 + 4 * elem
type IntArray []types.Int

func (*IntArray) Type() TagType             { return TagIntArray }
func (i IntArray) Size() int64              { return int64(4 + len(i)) }
func (*IntArray) Lookup(path string) Tagger { return nil }
func (i IntArray) String() string {
	var values []int32
	for _, elem := range i[:StringNum] {
		values = append(values, int32(elem))
	}

	s := "NBT_IntArray(size: %d) [% d, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(i), values, len(i)-StringNum)
}

func (i *IntArray) ReadFrom(reader io.Reader) (n int64, err error) {
	var length Int

	// Read length-prefix
	_, err = length.ReadFrom(reader)
	if err != nil {
		return
	}

	// Read length bytes
	arr := make([]types.Int, length.Int)
	for _, elem := range arr {
		_, err = elem.ReadFrom(reader)
		if err != nil {
			return
		}
	}

	*i = arr
	return
}

func (i *IntArray) WriteTo(writer io.Writer) (n int64, err error) {
	length := types.Int(len(*i))

	// Write length-prefix
	if _, err = length.WriteTo(writer); err != nil {
		return
	}

	// Then write int array
	for _, tag := range *i {
		if _, err = tag.WriteTo(writer); err != nil {
			return
		}
	}
	return
}
