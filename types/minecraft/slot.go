package minecraft

import (
	"bytes"
	"compress/gzip"
	"github.com/toqueteos/minero/proto/nbt"
	"io"

	// "github.com/toqueteos/minero/proto/nbt"
	"github.com/toqueteos/minero/util"
)

// http://wiki.vg/Slot_Data

// Slot
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
	var rw util.RWErrorHandler

	s.BlockId = rw.MustReadShort(r)
	if s.BlockId != -1 {
		s.Amount = byte(rw.MustReadByte(r))
		s.Damage = rw.MustReadShort(r)
		s.Length = rw.MustReadShort(r)
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
	var rw util.RWErrorHandler

	rw.MustWriteShort(w, s.BlockId)
	if s.BlockId != -1 {
		rw.MustWriteByte(w, int8(s.Amount))
		rw.MustWriteShort(w, s.Damage)
		rw.MustWriteShort(w, s.Length)
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
