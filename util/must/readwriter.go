package must

import (
	"bytes"
	"fmt"
	"io"

	"github.com/toqueteos/minero/types"
)

// ReadWriter keeps a counter of bytes read or write when dealing with
// 'io.Reader's and/or 'io.Writer's. Whenever an error is found all pending
// reads and writes are ignored.
type ReadWriter struct {
	Handler
	N int64
}

// Must provides a hook for foreign functions using 'io.Reader's or
// 'io.Writer's.
func (rw ReadWriter) Must(n int64, err error) {
	// Check iff there were no errors
	if rw.Err == nil {
		return
	}

	rw.N += n
	rw.Err = err
}

// Result returns all bytes read by r and the first error found.
func (rw ReadWriter) Result() (n int64, err error) {
	return rw.N, rw.Err
}

// Reset resets rw.
func (rw ReadWriter) Reset() {
	rw.N = 0
	rw.Err = nil
}

func (rw ReadWriter) ReadInt8(r io.Reader) (v int8) {
	if rw.Err != nil {
		return v
	}

	var t types.Int8
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadInt8: %v", err)
		return
	}
	rw.N += n

	return int8(t)
}

func (rw ReadWriter) ReadInt16(r io.Reader) (v int16) {
	if rw.Err != nil {
		return v
	}

	var t types.Int16
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadInt16: %v", err)
		return
	}
	rw.N += n

	return int16(t)
}

func (rw ReadWriter) ReadInt32(r io.Reader) (v int32) {
	if rw.Err != nil {
		return v
	}

	var t types.Int32
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadInt32: %v", err)
		return
	}
	rw.N += n

	return int32(t)
}

func (rw ReadWriter) ReadInt64(r io.Reader) (v int64) {
	if rw.Err != nil {
		return v
	}

	var t types.Int64
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadInt64: %v", err)
		return
	}
	rw.N += n

	return int64(t)
}

func (rw ReadWriter) ReadFloat32(r io.Reader) (v float32) {
	if rw.Err != nil {
		return v
	}

	var t types.Float32
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadFloat32: %v", err)
		return
	}
	rw.N += n

	return float32(t)
}

func (rw ReadWriter) ReadFloat64(r io.Reader) (v float64) {
	if rw.Err != nil {
		return v
	}

	var t types.Float64
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadFloat64: %v", err)
		return
	}
	rw.N += n

	return float64(t)
}

func (rw ReadWriter) ReadByteArray(r io.Reader, length int) (v []byte) {
	if rw.Err != nil {
		return v
	}

	var buf bytes.Buffer
	n, err := io.CopyN(&buf, r, int64(length))
	if err != nil {
		rw.Err = fmt.Errorf("ReadByteArray: %v", err)
		return
	}
	rw.N += n

	return buf.Bytes()
}

func (rw ReadWriter) WriteInt8(w io.Writer, value int8) {
	if rw.Err != nil {
		return
	}

	t := types.Int8(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteInt8: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteInt16(w io.Writer, value int16) {
	if rw.Err != nil {
		return
	}

	t := types.Int16(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteInt16: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteInt32(w io.Writer, value int32) {
	if rw.Err != nil {
		return
	}

	t := types.Int32(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteInt32: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteInt64(w io.Writer, value int64) {
	if rw.Err != nil {
		return
	}

	t := types.Int64(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteInt64: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteFloat32(w io.Writer, value float32) {
	if rw.Err != nil {
		return
	}

	t := types.Float32(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteFloat32: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteFloat64(w io.Writer, value float64) {
	if rw.Err != nil {
		return
	}

	t := types.Float64(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteFloat64: %v", err)
		return
	}
	rw.N += n
}

func (rw ReadWriter) WriteByteArray(w io.Writer, value []byte) {
	if rw.Err != nil {
		return
	}

	n, err := io.CopyN(w, bytes.NewBuffer(value), int64(len(value)))
	if err != nil {
		rw.Err = fmt.Errorf("WriteByteArray: %v", err)
		return
	}
	rw.N += n
}
