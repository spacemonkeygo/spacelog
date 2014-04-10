// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"fmt"
	"strconv"
	"strings"
)

type LogLevel int32

const (
	Debug    LogLevel = 10
	Info     LogLevel = 20
	Notice   LogLevel = 30
	Warning  LogLevel = 40
	Error    LogLevel = 50
	Critical LogLevel = 60
	// syslog has Alert
	// syslog has Emerg

	defaultLevel = Notice
)

func (l LogLevel) String() string {
	switch l.Match() {
	case Critical:
		return "CRIT"
	case Error:
		return "ERR"
	case Warning:
		return "WARN"
	case Notice:
		return "NOTE"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "UNSET"
	}
}

func (l LogLevel) Name() string {
	switch l.Match() {
	case Critical:
		return "critical"
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Notice:
		return "notice"
	case Info:
		return "info"
	case Debug:
		return "debug"
	default:
		return "unset"
	}
}

func (l LogLevel) Match() LogLevel {
	if l >= Critical {
		return Critical
	}
	if l >= Error {
		return Error
	}
	if l >= Warning {
		return Warning
	}
	if l >= Notice {
		return Notice
	}
	if l >= Info {
		return Info
	}
	if l >= Debug {
		return Debug
	}
	return 0
}

func LevelFromString(str string) (LogLevel, error) {
	switch strings.ToLower(str) {
	case "crit", "critical":
		return Critical, nil
	case "err", "error":
		return Error, nil
	case "warn", "warning":
		return Warning, nil
	case "note", "notice":
		return Notice, nil
	case "info":
		return Info, nil
	case "debug":
		return Debug, nil
	}
	val, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		return LogLevel(val), nil
	}
	return 0, fmt.Errorf("Invalid log level: %s", str)
}
