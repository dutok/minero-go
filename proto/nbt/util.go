package nbt

import (
	"io"
)

// Read reads an NBT compound from a byte stream. It doesn't handle compression.
func Read(r io.Reader) (c *Compound, err error) {
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

// Write writes an NBT compound to a byte stream. It doesn't handle compression.
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
