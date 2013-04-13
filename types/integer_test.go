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
		value := Int8(v)

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
		value := Int16(v)

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
		value := Int32(v)

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
		value := Int64(v)

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
