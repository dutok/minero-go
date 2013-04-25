package minecraft

import (
	"encoding/binary"
	"io"

	"github.com/toqueteos/minero/types"
)

type TypeEntry byte

const (
	TypeByte TypeEntry = iota
	TypeShort
	TypeInt
	TypeFloat
	TypeString
	TypeSlot
	TypeVector
)

func (t TypeEntry) New() (e Entry) {
	switch t {
	case 0:
		e = new(EntryByte)
	case 1:
		e = new(EntryShort)
	case 2:
		e = new(EntryInt)
	case 3:
		e = new(EntryFloat)
	case 4:
		e = new(EntryString)
	case 5:
		e = new(EntrySlot)
	case 6:
		e = new(EntryVector)
	}
	return
}

type Metadata struct {
	Entries map[byte]Entry
}

type Entry interface {
	io.ReaderFrom
	io.WriterTo
	Type() byte
	// Value() []byte
}

func NewMetadata() *Metadata {
	return &Metadata{make(map[byte]Entry)}
}

func (m *Metadata) ReadFrom(r io.Reader) (n int64, err error) {
	var (
		key byte
		nn  int64
	)

	for key != 0x7f {
		err = binary.Read(r, binary.BigEndian, &key)
		if err != nil {
			return 0, err
		}
		n += 1

		if key == 0x7f {
			break
		}

		typ := TypeEntry(key & 0xE0)
		index := key & 0x1F
		value := typ.New()

		nn, err = value.ReadFrom(r)
		if err != nil {
			return
		}
		n += nn

		m.Entries[index] = value
	}

	return
}

func (m *Metadata) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64
	for _, entry := range m.Entries {
		nn, err = entry.WriteTo(w)
		if err != nil {
			return n, err
		}
		n += nn
	}
	_, err = w.Write([]byte{0x7F})
	if err != nil {
		return n, err
	}
	n += 1

	return
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
	var nn int64
	nn, err = e.Data[0].ReadFrom(r)
	if err != nil {
		return
	}
	n += nn
	nn, err = e.Data[1].ReadFrom(r)
	if err != nil {
		return
	}
	n += nn
	nn, err = e.Data[2].ReadFrom(r)
	if err != nil {
		return
	}
	n += nn
	return
}

func (e *EntryVector) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64
	nn, err = e.Data[0].WriteTo(w)
	if err != nil {
		return
	}
	n += nn
	nn, err = e.Data[1].WriteTo(w)
	if err != nil {
		return
	}
	n += nn
	nn, err = e.Data[2].WriteTo(w)
	if err != nil {
		return
	}
	n += nn
	return
}
