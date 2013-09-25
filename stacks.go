// Copyright (C) 2013 Space Monkey, Inc.

package log

import (
    "flag"
    "log"
    "runtime"
)

var (
    stack_size = flag.Int("log.stack_trace_max_byte_length", 4096,
        "The max stack trace byte length to log")
)

func PrintWithStack(message string) {
    buf := make([]byte, *stack_size)
    buf = buf[:runtime.Stack(buf, false)]
    log.Printf("%s\n%s", message, buf)
}
