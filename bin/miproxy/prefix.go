package main

import (
	"io"
)

type PrefixWriter struct {
	rw     io.Writer
	Prefix string
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	var nn int
	prefix := []byte(pw.Prefix + " ")

	// Write prefix
	n, err = pw.rw.Write(prefix)
	if err != nil {
		return
	}
	nn += n

	// Write content
	n, err = pw.rw.Write(p)
	if err != nil {
		// Prefix write succesful
		n = nn
		return
	}
	nn += n

	return nn, nil
}
