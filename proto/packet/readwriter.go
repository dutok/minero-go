package packet

import (
	"fmt"
	"io"

	mct "github.com/minero/minero/types/minecraft"
	"github.com/minero/minero/util/must"
)

// MustReadWriter handles error handling while reading common types present in
// packets. It keeps a counter of bytes read and errors out on the first error.
type MustReadWriter struct {
	must.ReadWriter
}

func (rw MustReadWriter) ReadString(r io.Reader) (res string) {
	if rw.Err != nil {
		return
	}

	t := new(mct.String)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadString: %s", err)
		return
	}
	rw.N += n

	return string(*t)
}

func (rw MustReadWriter) WriteString(w io.Writer, value string) {
	if rw.Err != nil {
		return
	}

	t := mct.String(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteString: %s", err)
		return
	}
	rw.N += n
}

func (rw MustReadWriter) ReadBool(r io.Reader) (res bool) {
	if rw.Err != nil {
		return
	}

	t := new(mct.Bool)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadBool: %s", err)
		return
	}
	rw.N += n

	return bool(*t)
}

func (rw MustReadWriter) WriteBool(w io.Writer, value bool) {
	if rw.Err != nil {
		return
	}

	t := mct.Bool(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteBool: %s", err)
		return
	}
	rw.N += n
}

func (rw MustReadWriter) ReadSlot(r io.Reader) (res *mct.Slot) {
	if rw.Err != nil {
		return
	}

	t := mct.NewSlot()
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadSlot: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw MustReadWriter) WriteSlot(w io.Writer, value *mct.Slot) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteSlot: %s", err)
		return
	}
	rw.N += n
}

func (rw MustReadWriter) ReadObjectData(r io.Reader) (res *mct.ObjectData) {
	if rw.Err != nil {
		return
	}

	t := new(mct.ObjectData)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadObjectData: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw MustReadWriter) WriteObjectData(w io.Writer, value *mct.ObjectData) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteObjectData: %s", err)
		return
	}
	rw.N += n
}

func (rw MustReadWriter) ReadMetadata(r io.Reader) (res mct.Metadata) {
	if rw.Err != nil {
		return
	}

	t := mct.NewMetadata()
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("ReadMetadata: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw MustReadWriter) WriteMetadata(w io.Writer, value mct.Metadata) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("WriteMetadata: %s", err)
		return
	}
	rw.N += n
}
