// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iodebug implements an io debugging Reader and Writer. Based on testing/iotest.
package iodebug

import (
	"fmt"
	"io"
	"log"
	"strings"
)

func pretty(p []byte) string {
	var s = fmt.Sprintf("%02x", p)
	return strings.Replace(s, "00", ".", -1)

	// const pad = 32
	// var i = len(s)
	// var r []string
	// for i >= pad {
	// 	r = append(r, s[:pad])
	// 	s = s[pad:]
	// 	i = len(s)
	// }
	// if len(s) > 0 {
	// 	r = append(r, s)
	// }
	// return strings.Join(r, "\n")
}

type writeLogger struct {
	prefix string
	w      io.Writer
}

func (l *writeLogger) Write(p []byte) (n int, err error) {
	n, err = l.w.Write(p)
	if err != nil {
		log.Printf("%s %s: %v", l.prefix, pretty(p[:n]), err)
	} else {
		log.Printf("%s %s", l.prefix, pretty(p[:n]))
	}
	return
}

// NewWriteLogger returns a writer that behaves like w except
// that it logs (using log.Printf) each write to standard error,
// printing the prefix and the hexadecimal data written.
func NewWriteLogger(prefix string, w io.Writer) io.Writer {
	return &writeLogger{prefix, w}
}

type readLogger struct {
	prefix string
	r      io.Reader
}

func (l *readLogger) Read(p []byte) (n int, err error) {
	n, err = l.r.Read(p)
	if err != nil {
		log.Printf("%s %s: %v", l.prefix, pretty(p[:n]), err)
	} else {
		log.Printf("%s %s", l.prefix, pretty(p[:n]))
	}
	return
}

// NewReadLogger returns a reader that behaves like r except
// that it logs (using log.Print) each read to standard error,
// printing the prefix and the hexadecimal data written.
func NewReadLogger(prefix string, r io.Reader) io.Reader {
	return &readLogger{prefix, r}
}
