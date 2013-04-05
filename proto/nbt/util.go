package nbt

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

// Read reads an NBT compound from the given reader uncompressing contents from
// reader if they are gzip'd.
func Read(r io.Reader) (c *Compound, err error) {
	var rr io.Reader
	rr, _, err = GuessCompression(r)
	if err != nil {
		return nil, err
	}

	return ReadRaw(rr)
}

// ReadRaw reads an NBT compound from the given reader. It doesn't try to guess
// compression.
func ReadRaw(r io.Reader) (c *Compound, err error) {
	// Read TagType
	var tt TagType
	if _, err = tt.ReadFrom(r); err != nil {
		return nil, err
	}

	// TagType should be TagCompound
	if tt != TagCompound {
		return nil, ErrInvalidTop
	}

	// Read compound name
	var name String
	_, err = name.ReadFrom(r)
	if err != nil {
		return
	}

	// Read compound contents
	c = NewCompound(name.Value)
	_, err = c.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Write writes an NBT compound to the given writer. Doesn't handle compression.
func Write(w io.Writer, c *Compound) (err error) {
	if _, err = TagCompound.WriteTo(w); err != nil {
		return
	}

	nameTag := &String{c.Name}
	_, err = nameTag.WriteTo(w)
	if err != nil {
		return
	}

	_, err = c.WriteTo(w)
	if err != nil {
		return
	}

	return
}

// GuessCompression determines if a NBT io.Reader is gzip-compressed or not.
// Lots of inspiration here: http://goo.gl/pRNZl
func GuessCompression(r io.Reader) (rr io.Reader, gz bool, err error) {
	// It seems most (all?) gzip files contain a "magic number" prefix '0x1f'.
	const magicNum = 1

	var buf bytes.Buffer
	if nn, err := io.CopyN(&buf, r, magicNum); nn != magicNum || err != nil {
		return nil, false, err
	}

	// Check if reader has that prefix.
	if bytes.Equal(buf.Bytes(), []byte{0x1f}) {
		// File was gzip'd
		rr, err = gzip.NewReader(io.MultiReader(&buf, r))
		switch err {
		case gzip.ErrHeader:
			// File isn't gzip'd
		case nil:
			gz = true
			return
		default:
			log.Fatalln("nbt.GuessCompression:", err)
		}
	}

	// Concatenate whatever we read previously with all remaining contents.
	return io.MultiReader(&buf, r), false, nil
}
