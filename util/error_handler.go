package util

import (
	"bytes"
	"fmt"
	"io"

	"github.com/toqueteos/minero/types"
)

// RWErrorHandler errors out on the first error. It keeps a counter of bytes
// read.
type RWErrorHandler struct {
	N   int64
	Err error
}

// Must adds additional errors checks, allowing other functions or methods to
// make it fail sooner if any error is found.
func (rw *RWErrorHandler) Must(n int64, err error) {
	// Check iff there were no errors
	if rw.Err == nil {
		return
	}

	rw.N += n
	rw.Err = err
}

// Check adds additional errors checks, allowing other functions or methods to
// make it fail sooner if any error is found.
func (rw *RWErrorHandler) Check(err error) {
	// Check iff there were no errors
	if rw.Err == nil {
		rw.Err = err
	}
}

// Result returns all bytes read by r and the first error found.
func (rw *RWErrorHandler) Result() (n int64, err error) {
	return rw.N, rw.Err
}

// Reset resets rw.
func (rw *RWErrorHandler) Reset() {
	rw.N = 0
	rw.Err = nil
}

func (rw *RWErrorHandler) MustReadByte(r io.Reader) (v int8) {
	if rw.Err != nil {
		return v
	}

	var t types.Byte
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadByte: %v", err)
		return
	}
	rw.N += n

	return int8(t)
}

func (rw *RWErrorHandler) MustReadShort(r io.Reader) (v int16) {
	if rw.Err != nil {
		return v
	}

	var t types.Short
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadShort: %v", err)
		return
	}
	rw.N += n

	return int16(t)
}

func (rw *RWErrorHandler) MustReadInt(r io.Reader) (v int32) {
	if rw.Err != nil {
		return v
	}

	var t types.Int
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadInt: %v", err)
		return
	}
	rw.N += n

	return int32(t)
}

func (rw *RWErrorHandler) MustReadLong(r io.Reader) (v int64) {
	if rw.Err != nil {
		return v
	}

	var t types.Long
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadLong: %v", err)
		return
	}
	rw.N += n

	return int64(t)
}

func (rw *RWErrorHandler) MustReadFloat(r io.Reader) (v float32) {
	if rw.Err != nil {
		return v
	}

	var t types.Float
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadFloat: %v", err)
		return
	}
	rw.N += n

	return float32(t)
}

func (rw *RWErrorHandler) MustReadDouble(r io.Reader) (v float64) {
	if rw.Err != nil {
		return v
	}

	var t types.Double
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadDouble: %v", err)
		return
	}
	rw.N += n

	return float64(t)
}

func (rw *RWErrorHandler) MustReadByteArray(r io.Reader, length int) (v []byte) {
	if rw.Err != nil {
		return v
	}

	var buf bytes.Buffer
	n, err := io.CopyN(&buf, r, int64(length))
	if err != nil {
		rw.Err = fmt.Errorf("MustReadByteArray: %v", err)
		return
	}
	rw.N += n

	return buf.Bytes()
}

func (rw *RWErrorHandler) MustWriteByte(w io.Writer, value int8) {
	if rw.Err != nil {
		return
	}

	t := types.Byte(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteByte: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteShort(w io.Writer, value int16) {
	if rw.Err != nil {
		return
	}

	t := types.Short(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadShort: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteInt(w io.Writer, value int32) {
	if rw.Err != nil {
		return
	}

	t := types.Int(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteInt: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteLong(w io.Writer, value int64) {
	if rw.Err != nil {
		return
	}

	t := types.Long(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteLong: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteFloat(w io.Writer, value float32) {
	if rw.Err != nil {
		return
	}

	t := types.Float(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteFloat: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteDouble(w io.Writer, value float64) {
	if rw.Err != nil {
		return
	}

	t := types.Double(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteDouble: %v", err)
		return
	}
	rw.N += n
}

func (rw *RWErrorHandler) MustWriteByteArray(w io.Writer, value []byte) {
	if rw.Err != nil {
		return
	}

	n, err := io.CopyN(w, bytes.NewBuffer(value), int64(len(value)))
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteByteArray: %v", err)
		return
	}
	rw.N += n
}
