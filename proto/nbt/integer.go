package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Byte holds a single signed byte.
// TagType: 1, Size: 1 byte
type Byte struct{ Value int8 }

func (*Byte) Type() TagType          { return TagByte }
func (*Byte) Lookup(path string) Tag { return nil }

func (b *Byte) String() string {
	return fmt.Sprintf("NBT_Byte %d", b.Value)
}

func (b *Byte) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &b.Value)
}

func (b *Byte) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &b.Value)
}

// Short holds a single signed short.
// TagType: 2, Size: 2 bytes
type Short struct{ Value int16 }

func (*Short) Type() TagType          { return TagShort }
func (*Short) Lookup(path string) Tag { return nil }

func (s *Short) String() string {
	return fmt.Sprintf("NBT_Short %d", s.Value)
}

func (s *Short) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &s.Value)
}

func (s *Short) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &s.Value)
}

// Int holds a single signed integer.
// TagType: 3, Size: 4 bytes
type Int struct{ Value int32 }

func (*Int) Type() TagType          { return TagInt }
func (*Int) Lookup(path string) Tag { return nil }

func (i *Int) String() string {
	return fmt.Sprintf("NBT_Int %d", i.Value)
}

func (i *Int) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &i.Value)
}

func (i *Int) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &i.Value)
}

// Long holds a single signed long.
// TagType: 4, Size: 8 bytes
type Long struct{ Value int64 }

func (*Long) Type() TagType          { return TagLong }
func (*Long) Lookup(path string) Tag { return nil }

func (l *Long) String() string {
	return fmt.Sprintf("NBT_Long %d", l.Value)
}

func (l *Long) Read(reader io.Reader) (err error) {
	return binary.Read(reader, binary.BigEndian, &l.Value)
}

func (l *Long) Write(writer io.Writer) (err error) {
	return binary.Write(writer, binary.BigEndian, &l.Value)
}
