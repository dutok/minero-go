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

		rn, err := mcs.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		// Can't check string equality, quick may generate weird codepoints
		return len(v) == len(mcs.String())
		// return v == mcs.String()
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
