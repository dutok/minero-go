package minecraft

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/toqueteos/minero/proto/nbt"
	"github.com/toqueteos/minero/util/must"
)

// Slot
// http://wiki.vg/Slot_Data
type Slot struct {
	BlockId int16
	*InfoSlot
	Enchantments *nbt.Compound
}

type InfoSlot struct {
	Amount byte
	Damage int16
	Length int16
}

func (s *Slot) ReadFrom(r io.Reader) (n int64, err error) {
	var rw must.ReadWriter

	s.BlockId = rw.ReadInt16(r)
	if s.BlockId != -1 {
		s.Amount = byte(rw.ReadInt8(r))
		s.Damage = rw.ReadInt16(r)
		s.Length = rw.ReadInt16(r)
	}

	if s.Length != -1 {
		var buf bytes.Buffer
		rw.Must(io.CopyN(&buf, r, int64(s.Length)))

		gr, err := gzip.NewReader(&buf)
		rw.Check(err)
		s.Enchantments, err = nbt.Read(gr)
		rw.Check(err)
		rw.Check(gr.Close())
	}

	return rw.Result()
}

func (s *Slot) WriteTo(w io.Writer) (n int64, err error) {
	var rw must.ReadWriter

	rw.WriteInt16(w, s.BlockId)
	if s.BlockId != -1 {
		rw.WriteInt8(w, int8(s.Amount))
		rw.WriteInt16(w, s.Damage)
		rw.WriteInt16(w, s.Length)
	}

	// BUG(toqueteos): Write Compound as byte stream, then gzip it
	if s.Length != -1 {
		// var buf bytes.Buffer
		// nbt.Write(&buf, "", s.Enchantments)

		// gw := gzip.NewWriter(&buf)
		// rw.Must(io.Copy(w, gw))
		// rw.Check(gw.Close())
	}

	return rw.Result()
}
