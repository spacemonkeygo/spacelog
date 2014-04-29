// Copyright (C) 2014 Space Monkey, Inc.

/*
Package spacelog is a collection of interface lego bricks designed to help you
build a flexible logging system.

spacelog is loosely inspired by the Python logging library.

The basic interaction is between a Logger and a Handler. A Logger is
what the programmer typically interacts with for creating log messages. A
Logger will be at a given log level, and if log messages can clear that
specific logger's log level filter, they will be passed off to the Handler.

Loggers are instantiated from GetLogger and GetLoggerNamed.

A Handler is a very generic interface for handling log events. You can provide
your own Handler for doing structured JSON output or colorized output or
countless other things.

Provided are a simple TextHandler with a variety of log event templates and
TextOutput sinks, such as io.Writer, Syslog, and so forth.

Make sure to see the setup subpackage for easy and configurable logging setup
at process start
*/
package spacelog
