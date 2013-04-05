package nbt

import (
	"fmt"
	"io"
	"strings"
)

// Compounds hold a list of a named tags. Order is not guaranteed.
// TagType: 10, Size: 1 + 4 + elem * id_size bytes
type Compound struct {
	Name  string
	Value map[string]Tag
}

func NewCompound(name string) *Compound {
	c := &Compound{
		Name:  name,
		Value: make(map[string]Tag),
	}
	return c
}

func (c Compound) Type() TagType { return TagCompound }

func (c Compound) Size() (n int64) {
	// TagCompound + CompoundName + TagEnd
	n += 1 + 4 + 1
	for key, value := range c.Value {
		n += 1                   // TagType
		n += int64(4 + len(key)) // Key Name
		n += value.Size()        // Value
	}
	return
}

func (c Compound) Lookup(path string) Tag {
	components := strings.SplitN(path, "/", 2)
	tag, ok := c.Value[components[0]]
	if !ok {
		return nil
	}

	if len(components) >= 2 {
		return tag.Lookup(components[1])
	}

	return tag
}

func (c *Compound) String() string {
	var compound []string

	for k, v := range c.Value {
		compound = append(compound, fmt.Sprintf("%q: %v", k, v))
	}

	content := strings.Join(compound, ",\n")
	return fmt.Sprintf("NBT_Compound(size: %d) {\n%s\n}", len(c.Value), content)
}

// ReadFrom satifies io.ReaderFrom interface. TypeId is not decoded.
func (c *Compound) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Empty compound
	c.Value = make(map[string]Tag)

	for {
		// Read tag type
		var tt TagType
		nn, err = tt.ReadFrom(r)
		switch {
		case tt == TagEnd:
			return n + 1, nil // TagEnd is 1 byte
		case err != nil:
			return
		}
		n += nn

		// Read tag name
		var name String
		nn, err = name.ReadFrom(r)
		if err != nil {
			return
		}
		n += nn

		// Read payload
		var tag Tag
		// Corner case: Compounds keep a copy of their name
		if tt == TagCompound {
			tag = NewCompound(name.Value)
		} else {
			tag = tt.New()
		}
		if tag == nil {
			return n, fmt.Errorf("Compound.ReadFrom wrong TagType %d.", tt)
		}
		nn, err = tag.ReadFrom(r)
		if err != nil {
			return
		}
		n += nn

		// Save kv pair
		c.Value[name.Value] = tag
	}

	return
}

// WriteTo satifies io.WriterTo interface. TypeId is not encoded.
func (c *Compound) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64

	for name, tag := range c.Value {
		nn, err = tag.Type().WriteTo(w)
		if err != nil {
			return
		}
		n += nn

		nameTag := &String{name}
		nn, err = nameTag.WriteTo(w)
		if err != nil {
			return
		}
		n += nn

		nn, err = tag.WriteTo(w)
		if err != nil {
			return
		}
		n += nn
	}

	nn, err = TagEnd.New().WriteTo(w)
	if err != nil {
		return
	}
	n += nn

	return
}
