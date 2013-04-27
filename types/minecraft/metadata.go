package minecraft

import (
	"bytes"
	"github.com/toqueteos/minero/util/must"
	"io"

	"github.com/toqueteos/minero/types"
)

const (
	TypeByte   byte = iota // 0: 000
	TypeShort              // 1: 001
	TypeInt                // 2: 010
	TypeFloat              // 3: 011
	TypeString             // 4: 100
	TypeSlot               // 5: 101
	TypeVector             // 6: 110
)

func EntryFrom(e byte) (r Entry) {
	switch e {
	case 0:
		r = new(EntryByte)
	case 1:
		r = new(EntryShort)
	case 2:
		r = new(EntryInt)
	case 3:
		r = new(EntryFloat)
	case 4:
		r = new(EntryString)
	case 5:
		r = new(EntrySlot)
	case 6:
		r = new(EntryVector)
	}
	return
}

type Metadata struct {
	Entries map[byte]Entry
}

type Entry interface {
	Type() byte
	io.ReaderFrom
	io.WriterTo
}

func NewMetadata() Metadata {
	return Metadata{
		Entries: make(map[byte]Entry),
	}
}

func MetadataFrom(s []byte) (m Metadata, err error) {
	m = NewMetadata()
	_, err = m.ReadFrom(bytes.NewBuffer(s))
	return
}

func (m Metadata) ReadFrom(r io.Reader) (n int64, err error) {
	var rw must.ReadWriter

	var key byte
	for key != 0x7f {
		// Read type+key
		key = byte(rw.ReadInt8(r))
		if key == 0x7f {
			break
		}

		var (
			typ     byte  = key & 0xE0 >> 5
			index   byte  = key & 0x1F
			payload Entry = EntryFrom(typ)
		)

		// Read payload
		rw.Must(payload.ReadFrom(r))

		m.Entries[index] = payload
	}

	return rw.Result()
}

func (m Metadata) WriteTo(w io.Writer) (n int64, err error) {
	var rw must.ReadWriter
	var buf bytes.Buffer

	for index, payload := range m.Entries {
		buf.Reset()
		typ := payload.Type()
		rw.Check(buf.WriteByte(typ<<5 | (index & 0x1F)))

		// Write type+key & payload
		rw.Must(buf.WriteTo(w))
		rw.Must(payload.WriteTo(w))
	}

	buf.Reset()
	rw.Check(buf.WriteByte(0x7f))
	rw.Must(buf.WriteTo(w))

	return rw.Result()
}

type EntryByte struct{ types.Int8 }
type EntryShort struct{ types.Int16 }
type EntryInt struct{ types.Int32 }
type EntryFloat struct{ types.Float32 }
type EntryString struct{ String }
type EntrySlot struct{ Slot }
type EntryVector struct{ Data [3]types.Int32 }

func (e EntryByte) Type() byte   { return 0 }
func (e EntryShort) Type() byte  { return 1 }
func (e EntryInt) Type() byte    { return 2 }
func (e EntryFloat) Type() byte  { return 3 }
func (e EntryString) Type() byte { return 4 }
func (e EntrySlot) Type() byte   { return 5 }
func (e EntryVector) Type() byte { return 6 }

func (e *EntryVector) ReadFrom(r io.Reader) (n int64, err error) {
	var rw must.ReadWriter
	for i := 0; i < len(e.Data); i++ {
		rw.Must(e.Data[i].ReadFrom(r))
	}
	return rw.Result()
}

func (e *EntryVector) WriteTo(w io.Writer) (n int64, err error) {
	var rw must.ReadWriter
	for i := 0; i < len(e.Data); i++ {
		rw.Must(e.Data[i].WriteTo(w))
	}
	return rw.Result()
}
