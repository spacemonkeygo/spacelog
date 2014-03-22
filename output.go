// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"bytes"
	"io"
	"log"
	"log/syslog"
	"sync"
)

type TextOutput interface {
	Output(LogLevel, []byte)
}

type WriterOutput struct {
	w io.Writer
}

func NewWriterOutput(w io.Writer) *WriterOutput {
	return &WriterOutput{w: w}
}

func (o *WriterOutput) Output(_ LogLevel, message []byte) {
	o.w.Write(append(bytes.TrimRight(message, "\r\n"), '\n'))
}

type SyslogOutput struct {
	w *syslog.Writer
}

func NewSyslogOutput(facility syslog.Priority, tag string) (
	*SyslogOutput, error) {
	w, err := syslog.New(facility, tag)
	if err != nil {
		return nil, err
	}
	return &SyslogOutput{w: w}, nil
}

func (o *SyslogOutput) Output(level LogLevel, message []byte) {
	level = level.Match()
	for _, msg := range bytes.Split(message, []byte{'\n'}) {
		switch level {
		case Critical:
			o.w.Crit(string(msg))
		case Error:
			o.w.Err(string(msg))
		case Warning:
			o.w.Warning(string(msg))
		case Info:
			o.w.Info(string(msg))
		case Debug:
			fallthrough
		default:
			o.w.Debug(string(msg))
		}
	}
}

// StdlibOutput is for writing to the default logging system if no one has
// called Setup. If someone has called Setup though, this will most likely
// cause endless recursion.
type StdlibOutput struct{}

func (StdlibOutput) Output(_ LogLevel, message []byte) {
	log.Print(string(message))
}

type bufferMsg struct {
	level   LogLevel
	message []byte
}

type BufferedOutput struct {
	o       TextOutput
	c       chan bufferMsg
	running sync.Mutex
}

func NewBufferedOutput(output TextOutput, buffer int) *BufferedOutput {
	if buffer < 0 {
		buffer = 0
	}
	b := &BufferedOutput{
		o: output,
		c: make(chan bufferMsg, buffer)}
	go b.process()
	return b
}

func (b *BufferedOutput) Close() {
	close(b.c)
	b.running.Lock()
	b.running.Unlock()
}

func (b *BufferedOutput) Output(level LogLevel, message []byte) {
	b.c <- bufferMsg{level: level, message: message}
}

func (b *BufferedOutput) process() {
	b.running.Lock()
	defer b.running.Unlock()
	for {
		msg, open := <-b.c
		if !open {
			break
		}
		b.o.Output(msg.level, msg.message)
	}
}
