package types

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestByte(t *testing.T) {
	f := func(v int8) bool {
		var (
			b   = Byte(v)
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

		t.Logf("input signed byte: %q", v)
		t.Logf("output signed byte: %q", b)

		return v == int8(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestUByte(t *testing.T) {
	f := func(v byte) bool {
		var (
			b   = UByte(v)
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

		t.Logf("input unsigned byte: %q", v)
		t.Logf("output unsigned byte: %q", b)

		return v == byte(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestShort(t *testing.T) {
	f := func(v int16) bool {
		var (
			b   = Short(v)
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

		t.Logf("input short: %q", v)
		t.Logf("output short: %q", b)

		return v == int16(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestInt(t *testing.T) {
	f := func(v int32) bool {
		var (
			b   = Int(v)
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

		t.Logf("input integer: %q", v)
		t.Logf("output integer: %q", b)

		return v == int32(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	f := func(v int64) bool {
		var (
			b   = Long(v)
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

		t.Logf("input long: %q", v)
		t.Logf("output long: %q", b)

		return v == int64(b)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
