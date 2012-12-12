package nbt

import (
	"fmt"
	"io"
	"strings"
)

// Compounds hold a list of a named tags. Order is not guaranteed.
// TagType: 10, Size: 1 + 4 + elem * id_size bytes
type Compound struct{ Tags map[string]Tag }

func (*Compound) Type() TagType { return TagCompound }
func (c *Compound) Lookup(path string) (tag Tag) {
	var ok bool

	components := strings.SplitN(path, "/", 2)
	tag, ok = c.Tags[components[0]]
	if !ok {
		return nil
	}

	if len(components) >= 2 {
		return tag.Lookup(components[1])
	}

	return
}

func (c *Compound) String() string {
	var compound []string

	for k, v := range c.Tags {
		compound = append(compound, fmt.Sprintf("%q: %v", k, v))
	}

	content := strings.Join(compound, ", ")
	return fmt.Sprintf("NBT_Compound(size: %d) { %s }", len(c.Tags), content)
}

func (c *Compound) Read(reader io.Reader) (err error) {
	tags := make(map[string]Tag)
	var (
		tag  Tag
		name string
	)

	for {
		if name, tag, err = readNameTag(reader); err != nil {
			return
		}

		if tag == nil {
			break
		}

		tags[name] = tag
	}

	c.Tags = tags
	return
}

func (c *Compound) Write(writer io.Writer) (err error) {
	for name, tag := range c.Tags {
		err = writeNameTag(writer, name, tag)
		if err != nil {
			return
		}
	}

	// Write TagEnd
	err = TagEnd.Write(writer)

	return
}
