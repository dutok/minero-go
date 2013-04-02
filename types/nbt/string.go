package nbt

import (
	"encoding/binary"
	"io"
)

type String string

func (s *String) ReadFrom(reader io.Reader) (n int64, err error) {
	var length uint16

	// Here length-prefix is unsigned so we can't use Short
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return
	}
	n = 2

	// Read length bytes
	arr := make([]byte, length)
	_, err = io.ReadFull(reader, arr)
	if err != nil {
		return
	}
	n += int64(length)

	*s = String(arr)
	return
}

func (s *String) WriteTo(writer io.Writer) (n int64, err error) {
	length := uint16(len(*s))

	// Write unsigned length-prefix, we can't use Short
	err = binary.Write(writer, binary.BigEndian, &length)
	if err != nil {
		return
	}
	n = 2

	// Then write string bytes
	_, err = writer.Write([]byte(*s))
	if err != nil {
		return
	}
	n += int64(length)

	return
}
