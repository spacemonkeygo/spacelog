// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"path/filepath"
	"time"
)

type LogEvent struct {
	LoggerName string
	Level      LogLevel
	Message    string
	Filepath   string
	Line       int
	Timestamp  time.Time
}

func (LogEvent) Reset() string     { return "\x1b[0m" }
func (LogEvent) Bold() string      { return "\x1b[1m" }
func (LogEvent) Underline() string { return "\x1b[4m" }
func (LogEvent) Black() string     { return "\x1b[30m" }
func (LogEvent) Red() string       { return "\x1b[31m" }
func (LogEvent) Green() string     { return "\x1b[32m" }
func (LogEvent) Yellow() string    { return "\x1b[33m" }
func (LogEvent) Blue() string      { return "\x1b[34m" }
func (LogEvent) Magenta() string   { return "\x1b[35m" }
func (LogEvent) Cyan() string      { return "\x1b[36m" }
func (LogEvent) White() string     { return "\x1b[37m" }

func (l *LogEvent) Filename() string {
	if l.Filepath == "" {
		return ""
	}
	return filepath.Base(l.Filepath)
}

func (l *LogEvent) Time() string {
	return l.Timestamp.Format("15:04:05")
}

func (l *LogEvent) Date() string {
	return l.Timestamp.Format("2006/01/02")
}
