package types

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
		t.Logf("WriteTo: % x", buf.Bytes())
		t.Logf("WriteTo bytes written: %d", wn)

		rn, err := b.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}
		t.Logf("ReadFrom: % x", buf.Bytes())
		t.Logf("ReadFrom bytes read: %d", rn)

		t.Logf("input bool: %q", v)
		t.Logf("output bool: %q", b)

		return v == bool(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
