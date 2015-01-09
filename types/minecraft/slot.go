package minecraft

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/minero/minero/proto/nbt"
	"github.com/minero/minero/util/must"
)

// Slot
// http://wiki.vg/Slot_Data
type Slot struct {
	BlockId      int16
	Amount       byte
	Damage       int16
	Enchantments *nbt.Compound
}

func NewSlot() *Slot {
	return &Slot{
		BlockId:      -1,
		Amount:       1,
		Damage:       0,
		Enchantments: nil,
	}
}

func (s *Slot) ReadFrom(r io.Reader) (n int64, err error) {
	var rw must.ReadWriter

	s.BlockId = rw.ReadInt16(r)
	if s.BlockId == -1 {
		return rw.Result()
	}

	s.Amount = byte(rw.ReadInt8(r))
	s.Damage = rw.ReadInt16(r)
	Length := rw.ReadInt16(r)

	if Length == -1 {
		return rw.Result()
	}

	var br bytes.Buffer

	// Copy gzip'd NBT Compound
	gs := rw.ReadByteArray(r, int(Length))
	bn := bytes.NewBuffer(gs)

	// Ungzip byte array
	gr, err := gzip.NewReader(bn)
	rw.Check(err)
	rw.Must(io.Copy(&br, gr))
	rw.Check(gr.Close())

	// Read NBT Compound
	s.Enchantments, err = nbt.Read(&br)
	rw.Check(err)

	return rw.Result()
}

func (s *Slot) WriteTo(w io.Writer) (n int64, err error) {
	var rw must.ReadWriter

	rw.WriteInt16(w, s.BlockId)
	if s.BlockId == -1 {
		return rw.Result()
	}

	rw.WriteInt8(w, int8(s.Amount))
	rw.WriteInt16(w, s.Damage)

	if s.Enchantments == nil {
		rw.WriteInt16(w, -1)
		return rw.Result()
	}

	var bn, bw bytes.Buffer

	rw.Check(nbt.Write(&bn, s.Enchantments))
	gw := gzip.NewWriter(&bw)
	rw.Must(io.Copy(gw, &bn))
	rw.Check(gw.Close())

	rw.WriteInt16(w, int16(bw.Len()))
	rw.WriteByteArray(w, bw.Bytes())

	return rw.Result()
}
