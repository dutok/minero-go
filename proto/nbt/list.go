package nbt

import (
	"fmt"
	"io"
)

// List holds a list of nameless tags, all of the same type. The list is
// prefixed with the Type ID of the items it contains (1 byte), and the length
// of the list as a signed integer (4 bytes).
// TagType: 9, Size: 1 + 4 + elem * id_size bytes
type List struct {
	TagType
	Value []Tag
}

func (l List) Type() TagType          { return TagList }
func (l List) Lookup(path string) Tag { return nil }

func (l List) String() string {
	return fmt.Sprintf("NBT_List(size: %d) [ %s ]", len(l.Value), l.Value)
}

func (l *List) Read(reader io.Reader) (err error) {
	// Read TagType
	err = l.TagType.Read(reader)
	if err != nil {
		return
	}

	// Read length-prefix
	var length Int
	err = length.Read(reader)
	if err != nil {
		return
	}

	// Read list items
	list := make([]Tag, length.Value)
	for i, _ := range list {
		var tag Tag
		tag, err = l.TagType.New()
		if err != nil {
			return
		}
		err = tag.Read(reader)
		if err != nil {
			return
		}

		list[i] = tag
	}

	l.Value = list
	return
}

func (l *List) Write(writer io.Writer) (err error) {
	// Read TagType
	err = l.TagType.Write(writer)

	// Write TagType prefix
	tt := Byte{int8(l.TagType)}
	if err = tt.Write(writer); err != nil {
		return
	}

	length := Int{int32(len(l.Value))}
	if err = length.Write(writer); err != nil {
		return
	}

	for _, tag := range l.Value {
		if err = tag.Write(writer); err != nil {
			return
		}
	}

	return
}
