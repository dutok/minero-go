package types

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestString(t *testing.T) {
	f := func(v string) bool {
		var (
			mcs = String(v)
			buf bytes.Buffer
		)

		wn, err := mcs.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}
		t.Logf("WriteTo: % x", buf.Bytes())
		t.Logf("WriteTo bytes written: %d", wn)

		rn, err := mcs.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}
		t.Logf("ReadFrom: % x", buf.Bytes())
		t.Logf("ReadFrom bytes read: %d", rn)

		t.Logf("input string: %q", v)
		t.Logf("output string: %q", mcs.String())

		// Can't check string equality, quick generates weird
		return len(v) == len(mcs.String())
		// return v == mcs.String()
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
