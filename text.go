// Copyright (C) 2014 Space Monkey, Inc.

package spacelog

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
)

// TextHandler is the default implementation of the Handler interface. A
// TextHandler, on log events, makes LogEvent structures, passes them to the
// configured template, and then passes that output to a configured TextOutput
// interface.
type TextHandler struct {
	mtx      sync.RWMutex
	template *template.Template
	output   TextOutput
}

// NewTextHandler creates a Handler that creates LogEvents, passes them to
// the given template, and passes the result to output
func NewTextHandler(t *template.Template, output TextOutput) *TextHandler {
	return &TextHandler{template: t, output: output}
}

// Log makes a LogEvent, formats it with the configured template, then passes
// the output to configured output sink
func (h *TextHandler) Log(logger_name string, level LogLevel, msg string,
	calldepth int) {
	h.mtx.RLock()
	output, template := h.output, h.template
	h.mtx.RUnlock()
	event := LogEvent{
		LoggerName: logger_name,
		Level:      level,
		Message:    strings.TrimRight(msg, "\n\r"),
		Timestamp:  time.Now()}
	if calldepth >= 0 {
		_, event.Filepath, event.Line, _ = runtime.Caller(calldepth + 1)
	}
	var buf bytes.Buffer
	err := template.Execute(&buf, &event)
	if err != nil {
		output.Output(level, []byte(
			fmt.Sprintf("log format template failed: %s", err)))
		return
	}
	output.Output(level, buf.Bytes())
}

// SetTextTemplate changes the TextHandler's text formatting template
func (h *TextHandler) SetTextTemplate(t *template.Template) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.template = t
}

// SetTextOutput changes the TextHandler's TextOutput sink
func (h *TextHandler) SetTextOutput(output TextOutput) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.output = output
}
