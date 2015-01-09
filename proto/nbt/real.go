package nbt

import (
	"fmt"

	"github.com/minero/minero/types"
)

// Float holds a single IEEE-754 single-precision floating point number.
// TagType: 5, Size: 4 bytes
type Float32 struct {
	types.Float32
}

func (f Float32) Type() TagType          { return TagFloat }
func (f Float32) Size() int64            { return 4 }
func (f Float32) Lookup(path string) Tag { return nil }
func (f Float32) String() string         { return fmt.Sprintf("%f (float)", f) }

// Double holds a single IEEE-754 double-precision floating point number.
// TagType: 6, Size: 8 bytes
type Float64 struct {
	types.Float64
}

func (d Float64) Type() TagType          { return TagDouble }
func (d Float64) Size() int64            { return 8 }
func (d Float64) Lookup(path string) Tag { return nil }
func (d Float64) String() string         { return fmt.Sprintf("%f (double)", d) }
