package minecraft

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestString(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v string) bool {
		// Thanks to andrey mirtchovski for pointing out this:
		// https://groups.google.com/d/msg/golang-nuts/qeEvnU0yUr4/t5cMJCdPSNsJ
		// http://en.wikipedia.org/wiki/Mapping_of_Unicode_characters#Surrogates
		var rs []rune
		for _, r := range v {
			switch {
			case r >= 0xD800 && r <= 0xDBFF:
				continue
			case r >= 0xDC00 && r <= 0xDFFF:
				continue
			}
			rs = append(rs, r)
		}
		v = string(rs)

		value := String(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == string(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
