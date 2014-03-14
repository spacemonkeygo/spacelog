// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"sync"
	"sync/atomic"
)

type Logger struct {
	level      LogLevel
	name       string
	collection *LoggerCollection

	handler_mtx sync.RWMutex
	handler     Handler
}

func (l *Logger) Scope(name string) *Logger {
	return l.collection.getLogger(l.name+"."+name, l.getLevel(),
		l.getHandler())
}

func (l *Logger) setLevel(level LogLevel) {
	atomic.StoreInt32((*int32)(&l.level), int32(level))
}

func (l *Logger) getLevel() LogLevel {
	return LogLevel(atomic.LoadInt32((*int32)(&l.level)))
}

func (l *Logger) setHandler(handler Handler) {
	l.handler_mtx.Lock()
	defer l.handler_mtx.Unlock()
	l.handler = handler
}

func (l *Logger) getHandler() Handler {
	l.handler_mtx.RLock()
	defer l.handler_mtx.RUnlock()
	return l.handler
}
