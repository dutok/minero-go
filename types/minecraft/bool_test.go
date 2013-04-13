package minecraft

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestBool(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v bool) bool {
		value := Bool(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == bool(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
