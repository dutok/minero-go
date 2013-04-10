package minecraft

import (
	"github.com/toqueteos/minero/util"
	"io"
)

// ObjectData is a special data type for packet 0x17.
//
// Length and contents depend on the value of Data.
//
// Meaning of Data:
// - Item Frame (id 71); Orientation, 0~3: South, West, North, East
// - Falling Block (id 70); BlockType, BlockID | (Metadata << 0xC)
// - Projectiles; EntityId of thrower.
// - Splash Potions; PotionValue
type ObjectData struct {
	Data                   int32
	SpeedX, SpeedY, SpeedZ int16
}

func (o *ObjectData) ReadFrom(r io.Reader) (n int64, err error) {
	var rw util.RWErrorHandler

	o.Data = rw.MustReadInt(r)
	if o.Data != 0 {
		o.SpeedX = rw.MustReadShort(r)
		o.SpeedY = rw.MustReadShort(r)
		o.SpeedZ = rw.MustReadShort(r)
	}

	return rw.Result()
}

func (o *ObjectData) WriteTo(w io.Writer) (n int64, err error) {
	var rw util.RWErrorHandler

	rw.MustWriteInt(w, o.Data)
	if o.Data != 0 {
		rw.MustWriteShort(w, o.SpeedX)
		rw.MustWriteShort(w, o.SpeedY)
		rw.MustWriteShort(w, o.SpeedZ)
	}

	return rw.Result()
}
