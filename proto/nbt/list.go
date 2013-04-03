package nbt

import (
	"fmt"
	"io"

	"github.com/toqueteos/minero/types"
)

// List holds a list of nameless tags, all of the same type. The list is
// prefixed with the Type ID of the items it contains (1 byte), and the length
// of the list as a signed integer (4 bytes).
// TagType: 9, Size: 1 + 4 + elem * id_size bytes
type List struct {
	Typ   TagType
	Value []Tag
}

func (l List) Type() TagType { return TagList }
func (l List) Size() (n int64) {
	n = 5 // TagType + Name
	for _, elem := range l.Value {
		n += elem.Size()
	}
	return
}
func (l List) Lookup(path string) Tag { return nil }
func (l List) String() string {
	return fmt.Sprintf("NBT_List(size: %d) % s", len(l.Value), l.Value)
}

// ReadFrom satifies the io.ReaderFrom interface. Reads: list name, tag type and
// length elements and all list elements.
func (l *List) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read TagType
	nn, err = l.Typ.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length-prefix
	var length Int
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read list items
	if length.Int > 0 {
		l.Value = make([]Tag, length.Int)
		for _, elem := range l.Value {
			elem = l.Typ.New()
			if elem != nil {
				return
			}
			nn, err = elem.ReadFrom(r)
			if err != nil {
				return
			}
			n += nn
		}

	}

	return
}

func (l *List) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Read TagType
	nn, err = l.Typ.WriteTo(w)
	n += nn

	// Write TagType prefix
	tt := types.Byte(l.Typ)
	if nn, err = tt.WriteTo(w); err != nil {
		return
	}
	n += nn

	length := types.Int(len(l.Value))
	if nn, err = length.WriteTo(w); err != nil {
		return
	}
	n += nn

	for _, tag := range l.Value {
		if nn, err = tag.WriteTo(w); err != nil {
			return
		}
		n += nn
	}

	return
}
