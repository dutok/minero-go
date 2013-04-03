package nbt

import (
	"fmt"

	"github.com/toqueteos/minero/types"
)

// Byte holds a single signed byte.
// TagType: 1, Size: 1 byte
type Byte struct {
	types.Byte
}

func (b Byte) Type() TagType          { return TagByte }
func (b Byte) Size() int64            { return 1 }
func (b Byte) Lookup(path string) Tag { return nil }
func (b Byte) String() string         { return fmt.Sprintf("NBT_Byte %d", b.Byte) }

// Short holds a single signed short.
// TagType: 2, Size: 2 bytes
type Short struct {
	types.Short
}

func (s Short) Type() TagType          { return TagShort }
func (s Short) Size() int64            { return 2 }
func (s Short) Lookup(path string) Tag { return nil }
func (s Short) String() string         { return fmt.Sprintf("NBT_Short %d", s.Short) }

// Int holds a single signed integer.
// TagType: 3, Size: 4 bytes
type Int struct {
	types.Int
}

func (i Int) Type() TagType          { return TagInt }
func (i Int) Size() int64            { return 4 }
func (i Int) Lookup(path string) Tag { return nil }
func (i Int) String() string         { return fmt.Sprintf("NBT_Int %d", i.Int) }

// Long holds a single signed long.
// TagType: 4, Size: 8 bytes
type Long struct {
	types.Long
}

func (l Long) Type() TagType          { return TagLong }
func (l Long) Size() int64            { return 8 }
func (l Long) Lookup(path string) Tag { return nil }
func (l Long) String() string         { return fmt.Sprintf("NBT_Long %d", l.Long) }
