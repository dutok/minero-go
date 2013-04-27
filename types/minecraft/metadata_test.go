package minecraft

import (
	"bytes"
	"reflect"
	"testing"
)

var metadataTests = []struct {
	in  []byte
	out Metadata
}{
	{
		[]byte{0, 0, 0x21, 0x01, 0x2c, 6, 0, 127},
		Metadata{
			Entries: map[byte]Entry{
				0: &EntryByte{0},
				1: &EntryShort{300},
				6: &EntryByte{0},
			},
		},
	},
}

func TestMetadataRead(t *testing.T) {
	for index, tt := range metadataTests {
		out := NewMetadata()
		out.ReadFrom(bytes.NewBuffer(tt.in))

		if !reflect.DeepEqual(out, tt.out) {
			t.Fatalf("%d. Expecting %+v, got %+v.", index+1, tt.out, out)
		}
	}
}

func TestMetadataWrite(t *testing.T) {
	for index, tt := range metadataTests {
		var buf bytes.Buffer
		tt.out.WriteTo(&buf)

		out := NewMetadata()
		out.ReadFrom(&buf)

		if !reflect.DeepEqual(out, tt.out) {
			t.Fatalf("%d. Expecting %+v, got %+v.", index+1, tt.out, out)
		}
	}
}
