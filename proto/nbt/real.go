package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Float holds a single IEEE-754 single-precision floating point number.
// TagType: 5, Size: 4 bytes
type Float struct{ Value float32 }

func (*Float) Type() TagType          { return TagFloat }
func (*Float) Lookup(path string) Tag { return nil }

func (f *Float) String() string { return fmt.Sprintf("NBT_Float %f", f.Value) }

func (f *Float) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &f.Value)
}

func (f *Float) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &f.Value)
}

// Double holds a single IEEE-754 double-precision floating point number.
// TagType: 6, Size: 8 bytes
type Double struct{ Value float64 }

func (*Double) Type() TagType          { return TagDouble }
func (*Double) Lookup(path string) Tag { return nil }

func (d *Double) String() string { return fmt.Sprintf("NBT_Double %f", d.Value) }

func (d *Double) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &d.Value)
}

func (d *Double) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &d.Value)
}
