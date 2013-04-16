package nbt

import (
	"fmt"
	"io"
	"strings"

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
	var list []string

	for i, v := range l.Value {
		list = append(list, fmt.Sprintf("%d: %v,", i, v))
	}

	content := strings.Join(list, "")
	return fmt.Sprintf("List of %s with %d elems [%s]", l.Typ, len(l.Value), content)
}

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (l *List) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Read TagType
	nn, err = l.Typ.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read length-prefix
	var length Int32
	nn, err = length.ReadFrom(r)
	if err != nil {
		return
	}
	n += nn

	// Read list items
	if length.Int32 > 0 {
		l.Value = make([]Tag, length.Int32)
		for index, elem := range l.Value {
			elem = l.Typ.New()
			nn, err = elem.ReadFrom(r)
			if err != nil {
				return
			}
			l.Value[index] = elem
			n += nn
		}

	}

	return
}

// WriteTo satifies io.WriterTo interface. TypeId is not encoded.
func (l *List) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	// Write TagType prefix
	tt := types.Int8(l.Typ)
	if nn, err = tt.WriteTo(w); err != nil {
		return
	}
	n += nn

	length := types.Int32(len(l.Value))
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
