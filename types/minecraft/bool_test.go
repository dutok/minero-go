package minecraft

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestBool(t *testing.T) {
	f := func(v bool) bool {
		var (
			b   = Bool(v)
			buf bytes.Buffer
		)

		wn, err := b.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		rn, err := b.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == bool(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
