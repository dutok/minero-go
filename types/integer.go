package types

import (
	"encoding/binary"
	"io"
)

type Int8 int8

func (b *Int8) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, b)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (b *Int8) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, b)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// type UByte byte // Same implementation as Byte

type Int16 int16

func (s *Int16) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, s)
	if err != nil {
		return 0, err
	}
	return 2, nil
}

func (s *Int16) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return 0, err
	}
	return 2, nil
}

type Int32 int32

func (i *Int32) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, i)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (i *Int32) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, i)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type Int64 int64

func (l *Int64) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, l)
	if err != nil {
		return 0, err
	}
	return 8, nil
}

func (l *Int64) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, l)
	if err != nil {
		return 0, err
	}
	return 8, nil
}
