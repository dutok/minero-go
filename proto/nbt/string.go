package nbt

// http://play.golang.org/p/GSmvnu2jpR

import (
	"encoding/binary"
	"fmt"
	"io"
)

// String holds a length-prefixed UTF-8 string. The prefix is an unsigned short
// (2 bytes).
// TagType: 8, Size: 1 + (2 + elem) + (2 + elem)
type String struct {
	Value string
}

func (s String) Type() TagType          { return TagString }
func (s String) Size() int64            { return int64(2 + len(s.Value)) }
func (s String) Lookup(path string) Tag { return nil }
func (s String) String() string {
	return fmt.Sprintf("NBT_String(size: %d) %q", len(s.Value), s.Value)
}

func (s *String) ReadFrom(r io.Reader) (n int64, err error) {
	// unsigned short, can't use Short
	var length uint16
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return
	}

	if length > 0 {
		// Read length bytes
		arr := make([]byte, length)
		_, err = io.ReadFull(r, arr)
		if err != nil {
			return
		}
		n += int64(length)

		s.Value = string(arr)
	}

	return
}

func (s *String) WriteTo(w io.Writer) (n int64, err error) {
	// unsigned length-prefix, we can't use Short
	var length = uint16(len(s.Value))
	err = binary.Write(w, binary.BigEndian, &length)
	if err != nil {
		return
	}

	// Then write string bytes
	_, err = w.Write([]byte(s.Value))
	if err != nil {
		return
	}
	n += int64(length)

	return
}
