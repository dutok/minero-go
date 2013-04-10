package minecraft

import (
	"fmt"
	"io"

	mct "github.com/toqueteos/minero/types/minecraft"
	"github.com/toqueteos/minero/util"
)

// PacketErrorHandler errors out on the first error. It keeps a counter of bytes
// read.
type PacketErrorHandler struct {
	util.RWErrorHandler
}

func (rw *PacketErrorHandler) MustReadString(r io.Reader) (res string) {
	if rw.Err != nil {
		return
	}

	t := new(mct.String)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadString: %s", err)
		return
	}
	rw.N += n

	return string(*t)
}

func (rw *PacketErrorHandler) MustWriteString(w io.Writer, value string) {
	if rw.Err != nil {
		return
	}

	t := mct.String(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteString: %s", err)
		return
	}
	rw.N += n
}

func (rw *PacketErrorHandler) MustReadBool(r io.Reader) (res bool) {
	if rw.Err != nil {
		return
	}

	t := new(mct.Bool)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadBool: %s", err)
		return
	}
	rw.N += n

	return bool(*t)
}

func (rw *PacketErrorHandler) MustWriteBool(w io.Writer, value bool) {
	if rw.Err != nil {
		return
	}

	t := mct.Bool(value)
	n, err := t.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteBool: %s", err)
		return
	}
	rw.N += n
}

func (rw *PacketErrorHandler) MustReadSlot(r io.Reader) (res *mct.Slot) {
	if rw.Err != nil {
		return
	}

	t := new(mct.Slot)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadSlot: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw *PacketErrorHandler) MustWriteSlot(w io.Writer, value *mct.Slot) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteSlot: %s", err)
		return
	}
	rw.N += n
}

func (rw *PacketErrorHandler) MustReadObjectData(r io.Reader) (res *mct.ObjectData) {
	if rw.Err != nil {
		return
	}

	t := new(mct.ObjectData)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadObjectData: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw *PacketErrorHandler) MustWriteObjectData(w io.Writer, value *mct.ObjectData) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteObjectData: %s", err)
		return
	}
	rw.N += n
}

func (rw *PacketErrorHandler) MustReadMetadata(r io.Reader) (res *mct.Metadata) {
	if rw.Err != nil {
		return
	}

	t := new(mct.Metadata)
	n, err := t.ReadFrom(r)
	if err != nil {
		rw.Err = fmt.Errorf("MustReadMetadata: %s", err)
		return
	}
	rw.N += n

	return t
}

func (rw *PacketErrorHandler) MustWriteMetadata(w io.Writer, value *mct.Metadata) {
	if rw.Err != nil {
		return
	}

	n, err := value.WriteTo(w)
	if err != nil {
		rw.Err = fmt.Errorf("MustWriteMetadata: %s", err)
		return
	}
	rw.N += n
}
