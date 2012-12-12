package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

// String holds a length-prefixed UTF-8 string. The prefix is an unsigned short
// (2 bytes).
// TagType: 8, Size: 2 + elem * 1 bytes
type String struct{ Value string }

func (*String) Type() TagType          { return TagString }
func (*String) Lookup(path string) Tag { return nil }

func (s *String) String() string {
	return fmt.Sprintf("NBT_String(size: %d) %q", len(s.Value), s.Value)
}

func (s *String) Read(reader io.Reader) (err error) {
	var length uint16

	// Here length-prefix is unsigned so we can't use Short
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return
	}

	// Read length bytes
	arr := make([]byte, length)
	_, err = io.ReadFull(reader, arr)
	if err != nil {
		return
	}

	s.Value = string(arr)
	return
}

func (s *String) Write(writer io.Writer) (err error) {
	length := uint16(len(s.Value))

	// Write unsigned length-prefix, we can't use Short
	err = binary.Write(writer, binary.BigEndian, &length)
	if err != nil {
		return
	}

	// Then write string bytes
	_, err = writer.Write([]byte(s.Value))
	return
}
