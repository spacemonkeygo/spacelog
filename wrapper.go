// Copyright (C) 2013 Space Monkey, Inc.

package log

import (
	"bytes"
	"io"
)

type LogWrapper func(msg []byte) error

func (w LogWrapper) Write(msg []byte) (n int, err error) {
	n = len(msg)
	err = w(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

type newlineBreaker struct {
	w io.Writer
}

func NewlineBreaker(w io.Writer) io.Writer {
	return newlineBreaker{w: w}
}

func (w newlineBreaker) Write(msg []byte) (written int, err error) {
	for len(msg) > 0 {
		pos := bytes.Index(msg, []byte("\n"))
		if pos == -1 {
			n, err := w.w.Write(msg)
			return written + n, err
		}
		n, err := w.w.Write(msg[:pos+1])
		msg = msg[n:]
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}
