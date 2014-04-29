// Copyright (C) 2014 Space Monkey, Inc.

package spacelog

import (
	"path/filepath"
	"strings"
	"time"
)

// TermColors is a type that knows how to output terminal colors and formatting
type TermColors struct{}

// LogEvent is a type made by the default text handler for feeding to log
// templates. It has as much contextual data about the log event as possible.
type LogEvent struct {
	LoggerName string
	Level      LogLevel
	Message    string
	Filepath   string
	Line       int
	Timestamp  time.Time

	TermColors
}

// Reset resets the color palette for terminals that support color
func (TermColors) Reset() string     { return "\x1b[0m" }
func (TermColors) Bold() string      { return "\x1b[1m" }
func (TermColors) Underline() string { return "\x1b[4m" }
func (TermColors) Black() string     { return "\x1b[30m" }
func (TermColors) Red() string       { return "\x1b[31m" }
func (TermColors) Green() string     { return "\x1b[32m" }
func (TermColors) Yellow() string    { return "\x1b[33m" }
func (TermColors) Blue() string      { return "\x1b[34m" }
func (TermColors) Magenta() string   { return "\x1b[35m" }
func (TermColors) Cyan() string      { return "\x1b[36m" }
func (TermColors) White() string     { return "\x1b[37m" }

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

// LevelJustified returns the log level in string form justified so that all
// log levels take the same text width.
func (l *LogEvent) LevelJustified() (rv string) {
	rv = l.Level.String()
	if len(rv) < 5 {
		rv += strings.Repeat(" ", 5-len(rv))
	}
	return rv
}
