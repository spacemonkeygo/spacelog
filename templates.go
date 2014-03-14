// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"text/template"
)

func ColorizeLevel(level LogLevel) string {
	switch level.Match() {
	case Critical, Error:
		return LogEvent{}.Red()
	case Warning:
		return LogEvent{}.Magenta()
	case Info, Debug:
		return LogEvent{}.Green()
	}
	return ""
}

var (
	ColorTemplate = template.Must(template.New("color").Funcs(template.FuncMap{
		"ColorizeLevel": ColorizeLevel}).Parse(
		`{{.Blue}}{{.Date}} {{.Time}}{{.Reset}} ` +
			`{{.Bold}}{{ColorizeLevel .Level}}{{.Level}}{{.Reset}} ` +
			`{{.Underline}}{{.LoggerName}}{{.Reset}} ` +
			`{{if .Filename}}{{.Filename}}:{{.Line}} {{end}}- ` +
			`{{ColorizeLevel .Level}}{{.Message}}{{.Reset}}`))

	StandardTemplate = template.Must(template.New("standard").Parse(
		`{{.Date}} {{.Time}} ` +
			`{{.Level}} {{.LoggerName}} ` +
			`{{if .Filename}}{{.Filename}}:{{.Line}} {{end}}` +
			`- {{.Message}}`))

	SyslogTemplate = template.Must(template.New("syslog").Parse(
		`{{.Level}} {{.LoggerName}} ` +
			`{{if .Filename}}{{.Filename}}:{{.Line}} {{end}}` +
			`- {{.Message}}`))

	StdlibTemplate = template.Must(template.New("standard").Parse(
		`{{.Level}} {{.LoggerName}} ` +
			`{{if .Filename}}{{.Filename}}:{{.Line}} {{end}}` +
			`- {{.Message}}`))
)
