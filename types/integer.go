package types

import (
	"encoding/binary"
	"io"
)

type Byte int8

func (b *Byte) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, b)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (b *Byte) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, b)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// type UByte byte // Same implementation as Byte

type Short int16

func (s *Short) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, s)
	if err != nil {
		return 0, err
	}
	return 2, nil
}

func (s *Short) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return 0, err
	}
	return 2, nil
}

type Int int32

func (i *Int) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, i)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (i *Int) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, i)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type Long int64

func (l *Long) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, l)
	if err != nil {
		return 0, err
	}
	return 8, nil
}

func (l *Long) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, l)
	if err != nil {
		return 0, err
	}
	return 8, nil
}
