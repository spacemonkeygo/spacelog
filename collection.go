// Copyright (C) 2014 Space Monkey, Inc.

package spacelog

import (
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"
)

var (
	badChars = regexp.MustCompile("[^a-zA-Z0-9_.-]")
	slashes  = regexp.MustCompile("[/]")
)

func callerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown.unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown.unknown"
	}
	return badChars.ReplaceAllLiteralString(
		slashes.ReplaceAllLiteralString(
			strings.TrimPrefix(f.Name(), "code.spacemonkey.com/go/"), "."), "_")
}

type LoggerCollection struct {
	mtx     sync.Mutex
	loggers map[string]*Logger
	level   LogLevel
	handler Handler
}

func NewLoggerCollection() *LoggerCollection {
	return &LoggerCollection{
		loggers: make(map[string]*Logger),
		level:   DefaultLevel,
		handler: defaultHandler}
}

func (c *LoggerCollection) GetLogger() *Logger {
	return GetLoggerNamed(callerName())
}

func (c *LoggerCollection) getLogger(name string, level LogLevel,
	handler Handler) *Logger {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	logger, exists := c.loggers[name]
	if !exists {
		logger = &Logger{level: level,
			collection: c,
			name:       name,
			handler:    handler}
		c.loggers[name] = logger
	}
	return logger
}

func (c *LoggerCollection) GetLoggerNamed(name string) *Logger {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	logger, exists := c.loggers[name]
	if !exists {
		logger = &Logger{level: c.level,
			collection: c,
			name:       name,
			handler:    c.handler}
		c.loggers[name] = logger
	}
	return logger
}

func (c *LoggerCollection) SetLevel(re *regexp.Regexp, level LogLevel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if re == nil {
		c.level = level
	}
	for name, logger := range c.loggers {
		if re == nil || re.MatchString(name) {
			logger.setLevel(level)
		}
	}
}

func (c *LoggerCollection) SetHandler(re *regexp.Regexp, handler Handler) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if re == nil {
		c.handler = handler
	}
	for name, logger := range c.loggers {
		if re == nil || re.MatchString(name) {
			logger.setHandler(handler)
		}
	}
}

func (c *LoggerCollection) SetTextTemplate(re *regexp.Regexp,
	t *template.Template) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if re == nil {
		c.handler.SetTextTemplate(t)
	}
	for name, logger := range c.loggers {
		if re == nil || re.MatchString(name) {
			logger.getHandler().SetTextTemplate(t)
		}
	}
}

func (c *LoggerCollection) SetTextOutput(re *regexp.Regexp,
	output TextOutput) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if re == nil {
		c.handler.SetTextOutput(output)
	}
	for name, logger := range c.loggers {
		if re == nil || re.MatchString(name) {
			logger.getHandler().SetTextOutput(output)
		}
	}
}

var (
	DefaultLoggerCollection = NewLoggerCollection()
	GetLoggerNamed          = DefaultLoggerCollection.GetLoggerNamed
	SetLevel                = DefaultLoggerCollection.SetLevel
	SetHandler              = DefaultLoggerCollection.SetHandler
	SetTextTemplate         = DefaultLoggerCollection.SetTextTemplate
	SetTextOutput           = DefaultLoggerCollection.SetTextOutput
)

func GetLogger() *Logger {
	return GetLoggerNamed(callerName())
}
