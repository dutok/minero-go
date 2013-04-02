package nbt

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// Sets the maximum number of elements to show in string representations of
// types: NBT_ByteArray and NBT_IntArray.
const ArrayNum = 8

const (
	// Tag types. All these can be used to create a new tag, except TagEnd.
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

var (
	ErrEndTop     = errors.New("End tag found at top level.")
	ErrInvalidTop = errors.New("Expected compound at top level.")

	// String representation of each TagType
	tagName = map[TagType]string{
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
)

// TagType is the header byte value that identifies the type of tag(s). List &
// Compound types send TagType over the wire as a signed byte, using a int8 as
// underlying type allows us to assign TagType to Byte.
type TagType int8

func (tt *TagType) ReadFrom(r io.Reader) (n int64, err error) {
	// fmt.Printf("TagType was: %q.\n", tt)
	err = binary.Read(r, binary.BigEndian, tt)
	if err != nil {
		return 0, err
	}
	// fmt.Printf("TagType now is: %q.\n", tt)
	return 1, nil
}

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
		t = new(Byte)
	case TagShort:
		t = new(Short)
	case TagInt:
		t = new(Int)
	case TagLong:
		t = new(Long)
	case TagFloat:
		t = new(Float)
	case TagDouble:
		t = new(Double)
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

// Tag is the interface for all tags that can be represented in an NBT tree.
type Tag interface {
	io.ReaderFrom
	io.WriterTo
	// Name() string
	Type() TagType
	Size() int64
	Lookup(path string) Tag // Only Compound implements this
}

// Read reads an NBT compound from the given reader.
func Read(src io.Reader) (c *Compound, err error) {
	r, err := GuessCompression(src)
	if err != nil {
		return nil, err
	}

	// Read TagType
	var tt TagType
	if _, err = tt.ReadFrom(r); err != nil {
		return nil, err
	}

	// TagType should be TagCompound
	if tt != TagCompound {
		return nil, ErrInvalidTop
	}

	c = new(Compound)
	_, err = c.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Write writes an NBT compound to the given writer. Doesn't handle compression.
func Write(dst io.Writer, name string, tag *Compound) error {
	return nil
	// return writeNameTag(dst, name, tag)
}

// GuessCompression determines if a NBT io.Reader is gzip-compressed or not.
// Inspired on: http://goo.gl/pRNZl
func GuessCompression(r io.Reader) (rr io.Reader, err error) {
	// It seems most (all?) gzip files contain a "magic number" prefix 0x1f8b.
	const magicNumberRead = 2

	var buf bytes.Buffer
	if nn, err := io.CopyN(&buf, r, magicNumberRead); nn != 2 || err != nil {
		return nil, err
	}

	// Check if reader has that prefix.
	if !bytes.Equal(buf.Bytes(), []byte{0x1f, 0x8b}) {
		// Concatenate whatever we read previously with all remaining contents.
		return io.MultiReader(&buf, r), nil
	}

	// File was gzip'd
	rr, err = gzip.NewReader(io.MultiReader(&buf, r))
	switch err {
	case gzip.ErrHeader:
		// File isn't gzip'd
	case nil:
		return
	default:
		log.Fatalln("nbt.GuessCompression:", err)
	}

	return
}

// readNameTag reads tag type, name and tag contents from `src`. Useful for
// dealing with Compound structs.
// func readNameTag(r io.Reader) (name string, tag Tag, err error) {}

// writeNameTag writes tag type, name and tag contents to `w`. Useful for
// dealing with Compound structs.
// func writeNameTag(w io.Writer, name string, tag Tag) (err error) {
//  // Write tag type
//  _, err = tag.Type().WriteTo(w)
//  if err != nil {
//      return
//  }
//  // Write name
//  pathName := String{name}
//  _, err = pathName.WriteTo(w)
//  if err != nil {
//      return
//  }
//  // Write payload
//  _, err = tag.WriteTo(w)
//  return
// }
