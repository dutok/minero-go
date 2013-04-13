package nbt

import (
	"encoding/binary"
	"errors"
	"io"
)

// Sets the maximum number of elements to show in string representations of
// types: NBT_ByteArray and NBT_IntArray.
const ArrayNum = 8

var (
	ErrEndTop     = errors.New("End tag found at top level.")
	ErrInvalidTop = errors.New("Expected compound at top level.")
)

// Tag is the interface for all tags that can be represented in an NBT tree.
type Tag interface {
	io.ReaderFrom
	io.WriterTo
	Type() TagType
	Size() int64
	Lookup(path string) Tag
}

// TagType is the header byte value that identifies the type of tag(s). List &
// Compound types send TagType over the wire as a signed byte, using a int8 as
// underlying type allows us to assign TagType to Byte.
type TagType int8

const (
	// Tag types. All these can be used to create a new tag.
	TagEnd       TagType = iota // Size: 0
	TagByte                     // Size: 1
	TagShort                    // Size: 2
	TagInt                      // Size: 4
	TagLong                     // Size: 8
	TagFloat                    // Size: 4
	TagDouble                   // Size: 8
	TagByteArray                // Size: 4 + 1*elem
	TagString                   // Size: 2 + 4*elem
	TagList                     // Size: 1 + 4 + elem*len
	TagCompound                 // Size: varies
	TagIntArray                 // Size: 4 + 4*elem
)

// String representation of each TagType
var tagName = map[TagType]string{
	TagEnd:       "TagEnd",
	TagByte:      "TagByte",
	TagShort:     "TagShort",
	TagInt:       "TagInt",
	TagLong:      "TagLong",
	TagFloat:     "TagFloat",
	TagDouble:    "TagDouble",
	TagByteArray: "TagByteArray",
	TagString:    "TagString",
	TagList:      "TagList",
	TagCompound:  "TagCompound",
	TagIntArray:  "TagIntArray",
}

// ReadFrom satifies io.ReaderFrom interface.
func (tt *TagType) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, tt)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// WriteTo satifies io.WriterTo interface.
func (tt TagType) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, tt)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (tt TagType) String() string {
	if name, ok := tagName[tt]; ok {
		return name
	}
	return "TagErr"
}

func (tt TagType) New() (t Tag) {
	switch tt {
	case TagEnd:
		t = new(End)
	case TagByte:
		t = new(Int8)
	case TagShort:
		t = new(Int16)
	case TagInt:
		t = new(Int32)
	case TagLong:
		t = new(Int64)
	case TagFloat:
		t = new(Float32)
	case TagDouble:
		t = new(Float64)
	case TagByteArray:
		t = new(ByteArray)
	case TagString:
		t = new(String)
	case TagList:
		t = new(List)
	case TagCompound:
		t = new(Compound)
	case TagIntArray:
		t = new(IntArray)
	}
	return
}
