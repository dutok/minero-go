package types

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type String string

func (m *String) ReadFrom(r io.Reader) (n int64, err error) {
	var length uint16

	// Read length of string, 2 bytes
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return 0, err
	}

	// Read string, 2 bytes * length
	var contents = make([]uint16, length)
	err = binary.Read(r, binary.BigEndian, &contents)
	if err != nil {
		return 2, err
	}

	// Update String contents
	*m = String(readString(contents))

	return int64(2 + len(contents)*2), nil
}

func (m *String) WriteTo(w io.Writer) (n int64, err error) {
	var contents = writeString(string(*m))
	var length = uint16(len(contents))

	// Read size of string inside a short
	err = binary.Write(w, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.BigEndian, contents)
	if err != nil {
		return 2, err
	}

	return int64(2 + length*2), nil
}

// func (m *String) Read(p []byte) (n int, err error) {}
// func (m *String) Write(p []byte) (n int, err error) {}
func (m String) String() string {
	return string(m)
}

// writeString encodes a Go string to UCS-2 (UTF-16).
func writeString(s string) []uint16 {
	return utf16.Encode([]rune(s))
}

// readString decodes a UCS-2 (UTF-16) stream into a Go string.
func readString(u16 []uint16) string {
	return string(utf16.Decode(u16))
}
