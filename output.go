// Copyright (C) 2014 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spacelog

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

// WriterOutput is an io.Writer wrapper that matches the TextOutput interface
type WriterOutput struct {
	w io.Writer
}

// NewWriterOutput returns a TextOutput that writes messages to an io.Writer
func NewWriterOutput(w io.Writer) *WriterOutput {
	return &WriterOutput{w: w}
}

func (o *WriterOutput) Output(_ LogLevel, message []byte) {
	o.w.Write(append(bytes.TrimRight(message, "\r\n"), '\n'))
}

// SyslogOutput is a syslog client that matches the TextOutput interface
type SyslogOutput struct {
	w *syslog.Writer
}

// NewSyslogOutput returns a TextOutput object that writes to syslog using
// the given facility and tag. The log level will be determined by the log
// event.
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
		case Notice:
			o.w.Notice(string(msg))
		case Info:
			o.w.Info(string(msg))
		case Debug:
			fallthrough
		default:
			o.w.Debug(string(msg))
		}
	}
}

// StdlibOutput is a TextOutput that simply writes to the default Go stdlib
// logging system. It is the default. If you configure the Go stdlib to write
// to spacelog, make sure to provide a new TextOutput to your logging
// collection
type StdlibOutput struct{}

func (*StdlibOutput) Output(_ LogLevel, message []byte) {
	log.Print(string(message))
}

type bufferMsg struct {
	level   LogLevel
	message []byte
}

// BufferedOutput uses a channel to synchronize writes to a wrapped TextOutput
// and allows for buffering a limited amount of log events.
type BufferedOutput struct {
	o          TextOutput
	c          chan bufferMsg
	running    sync.Mutex
	close_once sync.Once
}

// NewBufferedOutput returns a BufferedOutput wrapping output with a buffer
// size of buffer.
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

// Close shuts down the BufferedOutput's processing
func (b *BufferedOutput) Close() {
	b.close_once.Do(func() {
		close(b.c)
	})
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
