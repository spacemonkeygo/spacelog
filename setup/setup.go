// Copyright (C) 2014 Space Monkey, Inc.

// Package setup provides simple helpers for configuring spacelog from flags.
package setup

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"math"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/SpaceMonkeyGo/spacelog"
)

var (
	output = flag.String("log.output", "stderr", "log output")
	level  = flag.String("log.level", spacelog.DefaultLevel.Name(),
		"base logger level")
	filter = flag.String("log.filter", "",
		"logger prefix to set level to the lowest level")
	format       = flag.String("log.format", "", "Format string to use")
	stdlog_level = flag.String("log.stdlevel", "warn",
		"logger level for stdlog integration")
	subproc = flag.String("log.subproc", "",
		"process to run for stdout/stderr-captured logging. The command is first "+
			"processed as a Go template that supports {{.Facility}}, {{.Level}}, "+
			"and {{.Name}} fields, and then passed to sh. If set, will redirect "+
			"stdout and stderr to the given process. A good default is "+
			"'setsid logger --priority {{.Facility}}.{{.Level}} --tag {{.Name}}'")
	buffer = flag.Int("log.buffer", 0, "the number of messages to buffer. "+
		"0 for no buffer")

	stdlog  = spacelog.GetLoggerNamed("stdlog")
	funcmap = template.FuncMap{"ColorizeLevel": spacelog.ColorizeLevel}
)

// SetFormatMethod adds functions to the template function map, such that
// command-line provided templates can call methods added to the map via this
// method. The map comes prepopulated with ColorizeLevel, but can be
// overridden. SetFormatMethod should be called (if at all) before one of
// this package's Setup methods.
func SetFormatMethod(name string, fn interface{}) {
	funcmap[name] = fn
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustSetup is the same as Setup, but panics instead of returning an error
func MustSetup(procname string) {
	must(Setup(procname))
}

// Setup is the same as SetupWithFacility, using the syslog.LOG_USER facility.
func Setup(procname string) error {
	return SetupWithFacility(procname, syslog.LOG_USER)
}

// MustSetupWithFacility is the same as SetupWithFacility, but panics instead
// of returning an error
func MustSetupWithFacility(procname string, facility syslog.Priority) {
	must(SetupWithFacility(procname, facility))
}

type subprocInfo struct {
	Facility string
	Level    string
	Name     string
}

// SetupWithFacility takes a given procname and a facility, and sets
// spacelog up with the given available flags and discovered configuration.
// SetupWithFacility supports (through flags):
//  * capturing stdout and stderr to a subprocess
//  * configuring the default level
//  * configuring log filters (enabling only some loggers)
//  * configuring the logging template
//  * configuring the output (a file, syslog, stdout, stderr)
//  * configuring log event buffering
//  * capturing all standard library logging with configurable log level
// It's pretty useless to call this method without parsing flags first, via
// flag.Parse() or flagfile.Load() or something.
func SetupWithFacility(procname string, facility syslog.Priority) error {
	if *subproc != "" {
		t, err := template.New("subproc").Parse(*subproc)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = t.Execute(&buf, &subprocInfo{
			Facility: fmt.Sprintf("%d", facility),
			Level:    fmt.Sprintf("%d", syslog.LOG_CRIT),
			Name:     procname})
		if err != nil {
			return err
		}
		err = spacelog.CaptureOutputToProcess("sh", "-c", string(buf.Bytes()))
		if err != nil {
			return err
		}
	}
	level_val, err := spacelog.LevelFromString(*level)
	if err != nil {
		return err
	}
	if level_val != spacelog.DefaultLevel {
		spacelog.SetLevel(nil, level_val)
	}
	if *filter != "" {
		re, err := regexp.Compile(*filter)
		if err != nil {
			return err
		}
		spacelog.SetLevel(re, spacelog.LogLevel(math.MinInt32))
	}
	var t *template.Template
	if *format != "" {
		var err error
		t, err = template.New("user").Funcs(funcmap).Parse(*format)
		if err != nil {
			return err
		}
	}
	var textout spacelog.TextOutput
	switch strings.ToLower(*output) {
	case "syslog":
		w, err := spacelog.NewSyslogOutput(facility, procname)
		if err != nil {
			return err
		}
		if t == nil {
			t = spacelog.SyslogTemplate
		}
		textout = w
	case "stdout":
		if t == nil {
			t = spacelog.ColorTemplate
		}
		textout = spacelog.NewWriterOutput(os.Stdout)
	case "stderr":
		if t == nil {
			t = spacelog.ColorTemplate
		}
		textout = spacelog.NewWriterOutput(os.Stderr)
	default:
		if t == nil {
			t = spacelog.StandardTemplate
		}
		fh, err := os.OpenFile(*output,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		textout = spacelog.NewWriterOutput(fh)
	}
	if *buffer > 0 {
		textout = spacelog.NewBufferedOutput(textout, *buffer)
	}
	spacelog.SetHandler(nil, spacelog.NewTextHandler(t, textout))
	log.SetFlags(log.Lshortfile)
	stdlog_level_val, err := spacelog.LevelFromString(*stdlog_level)
	if err != nil {
		return err
	}
	log.SetOutput(stdlog.WriterWithoutCaller(stdlog_level_val))
	return nil
}
