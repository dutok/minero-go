package nbt

import (
	"fmt"

	"github.com/toqueteos/minero/types"
)

// Float holds a single IEEE-754 single-precision floating point number.
// TagType: 5, Size: 4 bytes
type Float struct {
	types.Float
}

func (f Float) Type() TagType             { return TagFloat }
func (f Float) Size() int64               { return 4 }
func (f Float) Lookup(path string) Tagger { return nil }
func (f Float) String() string            { return fmt.Sprintf("NBT_Float %f", f) }

// Double holds a single IEEE-754 double-precision floating point number.
// TagType: 6, Size: 8 bytes
type Double struct {
	types.Double
}

func (d Double) Type() TagType             { return TagDouble }
func (d Double) Size() int64               { return 8 }
func (d Double) Lookup(path string) Tagger { return nil }
func (d Double) String() string            { return fmt.Sprintf("NBT_Double %f", d) }
