package nbt

// Package ideas.
//
// Instead of Read and Write methods implement equivalents from std:
// - WriteTo(io.Writer) (n int64, err error)  // io.WriterTo interface
// - ReadFrom(io.Reader) (n int64, err error) // io.ReaderFrom interface
//
// DRY away common behavior with a base type like: `type Common struct{}`

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// Sets the maximum number of elements to show in string representations of types: NBT_List, NBT_ByteArray and NBT_IntArray.
	StringNum = 10
)

var (
	ErrEndTop     = errors.New("End tag found at top level.")
	ErrInvalidTop = errors.New("Expected compound at top level.")
)

// TagType is the header byte value that identifies the type of tag(s). List &
// Compound types send TagType over the wire as a signed byte, using a int8 as
// underlying type allows us to assign TagType to Byte.
type TagType int8

func (tt *TagType) Read(reader io.Reader) (err error) {
	err = binary.Read(reader, binary.BigEndian, tt)
	return
}

func (tt TagType) Write(writer io.Writer) (err error) {
	err = binary.Write(writer, binary.BigEndian, tt)
	return
}

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

var tagProperties = map[TagType]string{
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

func (tt TagType) String() string {
	if name, ok := tagProperties[tt]; ok {
		return name
	}
	return "TagErr"
}

func (tt TagType) New() (tag Tag, err error) {
	switch tt {
	case TagByte:
		tag = new(Byte)
	case TagShort:
		tag = new(Short)
	case TagInt:
		tag = new(Int)
	case TagLong:
		tag = new(Long)
	case TagFloat:
		tag = new(Float)
	case TagDouble:
		tag = new(Double)
	case TagByteArray:
		tag = new(ByteArray)
	case TagString:
		tag = new(String)
	case TagList:
		tag = new(List)
	case TagCompound:
		tag = new(Compound)
	case TagIntArray:
		tag = new(IntArray)
	default:
		err = fmt.Errorf("Invalid NBT tag type '%d'.", tt)
	}
	return
}

// Tag is the interface for all tags that can be represented in an NBT tree.
type Tag interface {
	Type() TagType
	Read(io.Reader) error
	Write(io.Writer) error
	Lookup(path string) Tag // Only Compound implements this
}

// readNameTag reads tag type, name and tag contents from `src`. Useful for
// dealing with Compound structs.
func readNameTag(src io.Reader) (name string, tag Tag, err error) {
	// Read tag type
	var tt TagType
	if err = tt.Read(src); err != nil {
		return
	}

	// log.Printf("readNameTag read TagType: %q.\n", tt)

	if tt == TagEnd {
		// return name, tag, errors.New("TagEnd")
		return
	}

	// Read name
	var nameTag String
	if err = nameTag.Read(src); err != nil {
		return
	}
	name = nameTag.Value

	// Read tag
	tag, err = tt.New()
	if err != nil {
		return
	}
	if err = tag.Read(src); err != nil {
		return
	}

	return
}

// writeNameTag writes tag type, name and tag contents to `dst`. Useful for
// dealing with Compound structs.
func writeNameTag(dst io.Writer, name string, tag Tag) (err error) {
	// Write tag type
	err = tag.Type().Write(dst)
	if err != nil {
		return
	}

	// Write name
	var pathName = String{name}
	err = pathName.Write(dst)
	if err != nil {
		return
	}

	// Write tag
	err = tag.Write(dst)

	return
}

// Read reads an NBT compound from the given reader.
func Read(src io.Reader) (c *Compound, err error) {
	r := GuessCompression(src)

	name, tag, err := readNameTag(r)
	if err != nil {
		return nil, err
	}

	if name != "" {
		// 	return nil, errors.New("Root name should be empty.")
	}
	if tag == nil {
		return nil, ErrEndTop
	}

	c, ok := tag.(*Compound)
	if !ok {
		return nil, ErrInvalidTop
	}

	return c, nil
}

// Write writes an NBT compound to the given writer. Doesn't handle compression.
func Write(dst io.Writer, name string, tag *Compound) error {
	return writeNameTag(dst, name, tag)
}

// GuessCompression determines if a NBT io.Reader is compressed or not.
func GuessCompression(src io.Reader) (r io.Reader) {
	// Inspired on: http://goo.gl/pRNZl
	//
	//     " What I would do is give gzip.NewReader an io.Reader implementation
	//     which copies everything it Reads into a bytes.Buffer.
	//
	//     Then, if gzip.NewReader fails, make a new io.MultiReader
	//     concatenating that buffer of reads with "fh", which (potentially)
	//     data yet to be read. "

	var buf [3]bytes.Buffer
	var err error

	// Copy first bytes
	io.CopyN(io.MultiWriter(&buf[0], &buf[1], &buf[2]), src, 10)

	// Is it gzip'd?
	r, err = gzip.NewReader(io.MultiReader(&buf[0], src))
	if err == gzip.ErrHeader {
		// log.Println("It's not gzip compressed.")
	} else {
		// log.Println("It's gzip compressed!")
		return
	}

	// Is it zlib'd?
	r, err = zlib.NewReader(io.MultiReader(&buf[1], src))
	if err == zlib.ErrHeader {
		// log.Println("It's not zlib compressed.")
	} else {
		// log.Println("It's zlib compressed!")
		return
	}

	// log.Println("It's uncompressed!")

	// Concatenate whatever we read previously with all remaining contents.
	return io.MultiReader(&buf[2], src)
}
