package types

import (
	"encoding/binary"
	"io"
)

type Float32 float32

func (f *Float32) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, f)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (f *Float32) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type Float64 float64

func (d *Float64) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, d)
	if err != nil {
		return 0, err
	}
	return 8, nil
}

func (d *Float64) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, d)
	if err != nil {
		return 0, err
	}
	return 8, nil
}
