package minecraft

import (
	"io"

	"github.com/toqueteos/minero/util/must"
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
	var rw must.ReadWriter

	o.Data = rw.ReadInt32(r)
	if o.Data != 0 {
		o.SpeedX = rw.ReadInt16(r)
		o.SpeedY = rw.ReadInt16(r)
		o.SpeedZ = rw.ReadInt16(r)
	}

	return rw.Result()
}

func (o *ObjectData) WriteTo(w io.Writer) (n int64, err error) {
	var rw must.ReadWriter

	rw.WriteInt32(w, o.Data)
	if o.Data != 0 {
		rw.WriteInt16(w, o.SpeedX)
		rw.WriteInt16(w, o.SpeedY)
		rw.WriteInt16(w, o.SpeedZ)
	}

	return rw.Result()
}
