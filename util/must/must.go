package must

// Handler errors out on the first error. It keeps a counter of bytes
// read.
type Handler struct {
	Err error
}

// Check provides a hook for foreign functions returning errors.
func (h *Handler) Check(err error) {
	// Check iff there were no errors
	if h.Err == nil {
		h.Err = err
	}
}

// Result returns all bytes read by r and the first error found.
func (h *Handler) Result() error {
	return h.Err
}

// Reset resets h.
func (h *Handler) Reset() {
	h.Err = nil
}
