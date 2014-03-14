// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	output = flag.String("log.output", "stderr", "log output")
	level  = flag.String("log.level", "", "base logger level")
	filter = flag.String("log.filter", "",
		"logger prefix to set level to debug")
	format       = flag.String("log.format", "", "Format string to use")
	stdlog_level = flag.String("log.stdlevel", "info",
		"logger level for stdlog integration")
	syslog_binary = flag.String("log.subproc", "/usr/bin/logger",
		"process to run for stderr-captured logging")

	stdlog  = GetLoggerNamed("stdlog")
	funcmap = template.FuncMap{"ColorizeLevel": ColorizeLevel}
)

// SetFormatMethod should be called (if at all) before Setup
func SetFormatMethod(name string, fn interface{}) {
	funcmap[name] = fn
}

func Setup(procname string) error {
	return SetupWithFacility(procname, syslog.LOG_USER)
}

func SetupWithFacility(procname string, facility syslog.Priority) error {
	if *level != "" {
		level_val, err := LevelFromString(*level)
		if err != nil {
			return err
		}
		SetLevel(nil, level_val)
	}
	if *filter != "" {
		re, err := regexp.Compile(*filter)
		if err != nil {
			return err
		}
		SetLevel(re, Debug)
	}
	var t *template.Template
	if *format != "" {
		var err error
		t, err = template.New("user").Funcs(funcmap).Parse(*format)
		if err != nil {
			return err
		}
	}
	switch strings.ToLower(*output) {
	case "syslog":
		w, err := NewSyslogOutput(facility, procname)
		if err != nil {
			return err
		}
		if t == nil {
			t = SyslogTemplate
		}
		SetHandler(nil, NewTextHandler(t, w))
		if *syslog_binary != "" {
			err = CaptureOutputToProcess(procname, *syslog_binary)
			if err != nil {
				return err
			}
		}
	case "stdout":
		if t == nil {
			t = ColorTemplate
		}
		SetHandler(nil, NewTextHandler(t, NewWriterOutput(os.Stdout)))
	case "stderr":
		if t == nil {
			t = ColorTemplate
		}
		SetHandler(nil, NewTextHandler(t, NewWriterOutput(os.Stderr)))
	default:
		if t == nil {
			t = StandardTemplate
		}
		fh, err := os.OpenFile(*output,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		SetHandler(nil, NewTextHandler(t, NewWriterOutput(fh)))
	}
	log.SetFlags(log.Lshortfile)
	stdlog_level_val, err := LevelFromString(*stdlog_level)
	if err != nil {
		return err
	}
	log.SetOutput(stdlog.WriterWithoutCaller(stdlog_level_val))
	return nil
}
