// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Handler interface {
	// Log is called for every message. if calldepth is negative, caller
	// information is missing
	Log(logger_name string, level LogLevel, msg string, calldepth int)

	// these are expected to be no-ops on non-text-output handlers
	SetTextTemplate(t *template.Template)
	SetTextOutput(output TextOutput)
}

type TextHandler struct {
	mtx      sync.RWMutex
	template *template.Template
	output   TextOutput
}

// NewTextHandler creates a Handler that takes LogEvents, passes them to
// the given template, and passes the result to output
func NewTextHandler(t *template.Template, output TextOutput) *TextHandler {
	return &TextHandler{template: t, output: output}
}

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

func (h *TextHandler) SetTextTemplate(t *template.Template) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.template = t
}

func (h *TextHandler) SetTextOutput(output TextOutput) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	h.output = output
}

var (
	defaultHandler = NewTextHandler(ColorTemplate,
		NewWriterOutput(os.Stderr))
)
