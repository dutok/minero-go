package types

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestFloat(t *testing.T) {
	f := func(v float32) bool {
		var (
			b   = Float(v)
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

		t.Logf("input float: %q", v)
		t.Logf("output float: %q", b)

		return v == float32(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestDouble(t *testing.T) {
	f := func(v float64) bool {
		var (
			b   = Double(v)
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

		t.Logf("input double: %q", v)
		t.Logf("output double: %q", b)

		return v == float64(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
