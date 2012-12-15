package types

import (
	"encoding/binary"
	"io"
)

type Float float32

func (f *Float) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, f)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (f *Float) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type Double float64

func (d *Double) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, d)
	if err != nil {
		return 0, err
	}
	return 8, nil
}

func (d *Double) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, d)
	if err != nil {
		return 0, err
	}
	return 8, nil
}
