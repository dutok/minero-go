package types

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestFloat(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v float32) bool {
		value := Float32(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == float32(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestDouble(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v float64) bool {
		value := Float64(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == float64(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
