package nbt

import (
	"fmt"

	"github.com/minero/minero/types"
)

// Byte holds a single signed byte.
// TagType: 1, Size: 1 byte
type Int8 struct {
	types.Int8
}

func (b Int8) Type() TagType          { return TagByte }
func (b Int8) Size() int64            { return 1 }
func (b Int8) Lookup(path string) Tag { return nil }
func (b Int8) String() string         { return fmt.Sprintf("%d (byte)", b.Int8) }

// Short holds a single signed short.
// TagType: 2, Size: 2 bytes
type Int16 struct {
	types.Int16
}

func (s Int16) Type() TagType          { return TagShort }
func (s Int16) Size() int64            { return 2 }
func (s Int16) Lookup(path string) Tag { return nil }
func (s Int16) String() string         { return fmt.Sprintf("%d (short)", s.Int16) }

// Int holds a single signed integer.
// TagType: 3, Size: 4 bytes
type Int32 struct {
	types.Int32
}

func (i Int32) Type() TagType          { return TagInt }
func (i Int32) Size() int64            { return 4 }
func (i Int32) Lookup(path string) Tag { return nil }
func (i Int32) String() string         { return fmt.Sprintf("%d (int)", i.Int32) }

// Long holds a single signed long.
// TagType: 4, Size: 8 bytes
type Int64 struct {
	types.Int64
}

func (l Int64) Type() TagType          { return TagLong }
func (l Int64) Size() int64            { return 8 }
func (l Int64) Lookup(path string) Tag { return nil }
func (l Int64) String() string         { return fmt.Sprintf("%d (long)", l.Int64) }
