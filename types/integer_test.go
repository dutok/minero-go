package types

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestByte(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v int8) bool {
		value := Byte(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == int8(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestShort(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v int16) bool {
		value := Short(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == int16(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestInt(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v int32) bool {
		value := Int(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == int32(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	f := func(v int64) bool {
		value := Long(v)

		_, err = value.WriteTo(&buf)
		if err != nil {
			t.Error(err)
		}

		_, err = value.ReadFrom(&buf)
		if err != nil {
			t.Error(err)
		}

		return v == int64(value)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
