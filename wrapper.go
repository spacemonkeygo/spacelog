// Copyright (C) 2013 Space Monkey, Inc.

package log

type LogWrapper func(msg []byte) error

func (w LogWrapper) Write(msg []byte) (n int, err error) {
    n = len(msg)
    err = w(msg)
    if err != nil {
        return 0, err
    }
    return n, nil
}
