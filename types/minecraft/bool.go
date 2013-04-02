package types

import (
	"encoding/binary"
	"io"
)

type Bool bool

func (b *Bool) ReadFrom(r io.Reader) (n int64, err error) {
	var btemp byte

	err = binary.Read(r, binary.BigEndian, &btemp)
	if err != nil {
		return 0, err
	}

	// Update stored value
	*b = btemp == 0x01

	return 1, nil
}

func (b *Bool) WriteTo(w io.Writer) (n int64, err error) {
	var btemp byte

	if *b == true {
		btemp = 0x01
	}

	err = binary.Write(w, binary.BigEndian, btemp)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
