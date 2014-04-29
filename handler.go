// Copyright (C) 2014 Space Monkey, Inc.

package spacelog

import (
	"text/template"
)

// Handler is an interface that knows how to process log events. This is the
// basic interface type for building a logging system. If you want to route
// structured log data somewhere, you would implement this interface.
type Handler interface {
	// Log is called for every message. if calldepth is negative, caller
	// information is missing
	Log(logger_name string, level LogLevel, msg string, calldepth int)

	// These two calls are expected to be no-ops on non-text-output handlers
	SetTextTemplate(t *template.Template)
	SetTextOutput(output TextOutput)
}

// HandlerFunc is a type to make implementation of the Handler interface easier
type HandlerFunc func(logger_name string, level LogLevel, msg string,
	calldepth int)

// Log simply calls f(logger_name, level, msg, calldepth)
func (f HandlerFunc) Log(logger_name string, level LogLevel, msg string,
	calldepth int) {
	f(logger_name, level, msg, calldepth)
}

// SetTextTemplate is a no-op
func (HandlerFunc) SetTextTemplate(t *template.Template) {}

// SetTextOutput is a no-op
func (HandlerFunc) SetTextOutput(output TextOutput) {}

var (
	defaultHandler = NewTextHandler(StdlibTemplate,
		&StdlibOutput{})
)
