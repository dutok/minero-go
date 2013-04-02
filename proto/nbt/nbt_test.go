package nbt

import (
	"bytes"
	"reflect"
	"testing"
)

var simpleTests = []struct {
	name string
	in   []byte
	out  *Compound
	err  error
}{
	{
		"invalidTop",
		// Root tag isn't Compound
		[]byte{
			0x01, 0x00, 0x04, 0x61, 0x62, 0x63, 0x64, 0x11,
		},
		nil,
		ErrInvalidTop,
	},
	{
		"compound byte",
		[]byte{
			0x0a, 0x00, 0x00, 0x01, 0x00, 0x04, 0x61, 0x62, 0x63, 0x64,
			0x11, 0x00,
		},
		&Compound{map[string]Tag{"abcd": &Byte{0x11}}},
		nil,
	},
	{
		"compound short",
		[]byte{
			0x0a, 0x00, 0x00, 0x02, 0x00, 0x04, 0x61, 0x62, 0x63, 0x64,
			0x11, 0x22, 0x00,
		},
		&Compound{map[string]Tag{"abcd": &Short{0x1122}}},
		nil,
	},
	{
		"compound int",
		[]byte{
			0x0a, 0x00, 0x00, 0x03, 0x00, 0x04, 0x61, 0x62, 0x63, 0x64,
			0x11, 0x22, 0x33, 0x44, 0x00,
		},
		&Compound{map[string]Tag{"abcd": &Int{0x11223344}}},
		nil,
	},
}

func TestRead(t *testing.T) {
	for _, test := range simpleTests {
		// Read byte blob
		r := bytes.NewBuffer(test.in)
		c, err := Read(r)

		// Expected compound?
		if !reflect.DeepEqual(c, test.out) {
			t.Logf("Compound: %v", c)
			t.Logf("Expected: %v", test.out)
			t.Fatalf("%s: Compounds don't match.", test.name)
		}

		// Expected errors?
		if test.err != err {
			t.Fatalf("%s: Errors don't match.", test.name)
		}
	}
}
