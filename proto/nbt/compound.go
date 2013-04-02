package nbt

import (
	"fmt"
	"io"
	"strings"
)

// Compounds hold a list of a named tags. Order is not guaranteed.
// TagType: 10, Size: 1 + 4 + elem * id_size bytes
type Compound struct {
	Value map[string]Tag
}

func NewCompound(name string) *Compound {
	c := &Compound{
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

	return
}

func (c *Compound) String() string {
	var compound []string

	for k, v := range c.Value {
		compound = append(compound, fmt.Sprintf("%q: %v", k, v))
	}

	content := strings.Join(compound, ",\n")
	return fmt.Sprintf("NBT_Compound(size: %d) {\n%s\n}", len(c.Value), content)
}

func (c *Compound) ReadFrom(r io.Reader) (n int64, err error) {
	var nn int64

	// Empty compound
	c.Value = make(map[string]Tag)

	// Top level compound name is read but ignored
	var name String
	nn, err = name.ReadFrom(r)
	if err != nil {
		fmt.Println("Compound.ReadFrom name error:", err)
		return
	}
	fmt.Printf("Compound name: %q\n", name.Value)
	n += nn

	var tt TagType
	for {
		// Read tag type
		nn, err = tt.ReadFrom(r)
		// End of compound?
		// if tt == TagEnd || err == io.EOF {
		if tt == TagEnd {
			fmt.Println("Compound.ReadFrom TagEnd/EOF found error:", err)
			return n + 1, nil // TagEnd is 1 byte
		}
		fmt.Println("Compound.ReadFrom TagType:", tt)
		if err != nil {
			fmt.Println("Compound.ReadFrom TagType error:", err)
			return
		}
		n += nn

		// Read tag name
		var name String
		nn, err = name.ReadFrom(r)
		if err != nil {
			fmt.Println("Compound.ReadFrom tag name error:", err)
			return
		}
		fmt.Println("Compound.ReadFrom tag name:", err)
		n += nn

		// Read payload
		tag := tt.New()
		if tag == nil {
			return n, fmt.Errorf("Compound.ReadFrom wrong TagType %d.", tt)
		}
		nn, err = tag.ReadFrom(r)
		if err != nil {
			fmt.Println("Compound.ReadFrom payload error:", err)
			return
		}
		n += nn

		// Save kv pair
		c.Value[name.Value] = tag

		fmt.Printf("%s payload: %q\n", tt, tag)
	}

	return
}

// TODO(toqueteos): Incomplete
func (c *Compound) WriteTo(w io.Writer) (n int64, err error) {
	// for name, tag := range c.Value {
	// 	err = writeNameTag(w, name, tag)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	// // Write TagEnd
	// TagEnd.New().WriteTo(w)
	return 0, nil
}
