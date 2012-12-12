package nbt

import (
	"fmt"
	"io"
)

// ByteArray holds a length-prefixed array of signed bytes. The prefix is a
// signed integer (4 bytes).
// TagType: 7, Size: 4 + elem * 1 bytes
type ByteArray struct{ Value []byte }

func (*ByteArray) Type() TagType          { return TagByteArray }
func (*ByteArray) Lookup(path string) Tag { return nil }

func (b *ByteArray) String() string {
	s := "NBT_ByteArray(size: %d) [% x, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(b.Value), b.Value[:StringNum], len(b.Value)-StringNum)
}

func (b *ByteArray) Read(reader io.Reader) (err error) {
	var length Int

	// Read length-prefix
	err = length.Read(reader)
	if err != nil {
		return
	}

	// Read length bytes
	arr := make([]byte, length.Value)
	_, err = io.ReadFull(reader, arr)
	if err != nil {
		return
	}

	b.Value = arr
	return
}

func (b *ByteArray) Write(writer io.Writer) (err error) {
	length := Int{int32(len(b.Value))}

	// Write length-prefix
	if err = length.Write(writer); err != nil {
		return
	}

	// Then write byte array
	_, err = writer.Write(b.Value)
	return
}

// IntArray holds a length-prefixed array of signed integers. The prefix is a
// signed integer (4 bytes) and indicates the number of 4 byte integers.
// TagType: 11, Size: 4 + 4 * elem
type IntArray struct {
	Value []Int
}

func (*IntArray) Type() TagType          { return TagIntArray }
func (*IntArray) Lookup(path string) Tag { return nil }

func (i *IntArray) String() string {
	var values []int32
	for _, elem := range i.Value[:StringNum] {
		values = append(values, elem.Value)
	}

	s := "NBT_IntArray(size: %d) [% d, ... (%d elem(s) more)]"
	return fmt.Sprintf(s, len(i.Value), values, len(i.Value)-StringNum)
}

func (i *IntArray) Read(reader io.Reader) (err error) {
	var length Int

	// Read length-prefix
	err = length.Read(reader)
	if err != nil {
		return
	}

	// Read length bytes
	arr := make([]Int, length.Value)
	for _, elem := range arr {
		err = elem.Read(reader)
		if err != nil {
			return
		}
	}

	i.Value = arr
	return
}

func (i *IntArray) Write(writer io.Writer) (err error) {
	length := Int{int32(len(i.Value))}

	// Write length-prefix
	if err = length.Write(writer); err != nil {
		return
	}

	// Then write int array
	for _, tag := range i.Value {
		if err = tag.Write(writer); err != nil {
			return
		}
	}
	return
}
