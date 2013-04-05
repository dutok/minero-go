package nbt

import (
	"encoding/binary"
	"io"
)

// End marks end of a Compound.
// TagType: 0, Size: 1 byte
type End struct{}

func (e End) Name() string           { return "TagEnd" }
func (e End) Type() TagType          { return TagEnd }
func (e End) Size() int64            { return 1 }
func (e End) Lookup(path string) Tag { return nil }
func (e End) String() string         { return "NBT_End" }

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (e *End) ReadFrom(r io.Reader) (n int64, err error) {
	var tid byte

	err = binary.Read(r, binary.BigEndian, &tid)
	if err != nil {
		return
	}
	n = 1
	return
}

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (e *End) WriteTo(w io.Writer) (n int64, err error) {
	var tid byte

	err = binary.Write(w, binary.BigEndian, &tid)
	if err != nil {
		return
	}
	n = 1
	return
}
