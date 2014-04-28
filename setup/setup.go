// Copyright (C) 2014 Space Monkey, Inc.

package setup

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"math"
	"os"
	"regexp"
	"strings"
	"text/template"

	space_log "code.spacemonkey.com/go/space/log"
)

var (
	output = flag.String("log.output", "stderr", "log output")
	level  = flag.String("log.level", space_log.DefaultLevel.Name(),
		"base logger level")
	filter = flag.String("log.filter", "",
		"logger prefix to set level to the lowest level")
	format       = flag.String("log.format", "", "Format string to use")
	stdlog_level = flag.String("log.stdlevel", "warn",
		"logger level for stdlog integration")
	subproc = flag.String("log.subproc", "",
		"process to run for stdout/stderr-captured logging. If set (usually to "+
			"/usr/bin/logger), will redirect stdout and stderr to the given "+
			"process. process should take --priority <facility>.<level> and "+
			"--tag <name> options")
	buffer = flag.Int("log.buffer", 0, "the number of messages to buffer. "+
		"0 for no buffer")

	stdlog  = space_log.GetLoggerNamed("stdlog")
	funcmap = template.FuncMap{"ColorizeLevel": space_log.ColorizeLevel}
)

// SetFormatMethod should be called (if at all) before Setup
func SetFormatMethod(name string, fn interface{}) {
	funcmap[name] = fn
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustSetup(procname string) {
	must(Setup(procname))
}

func Setup(procname string) error {
	return SetupWithFacility(procname, syslog.LOG_USER)
}

func MustSetupWithFacility(procname string, facility syslog.Priority) {
	must(SetupWithFacility(procname, facility))
}

func SetupWithFacility(procname string, facility syslog.Priority) error {
	if *subproc != "" {
		err := space_log.CaptureOutputToProcess("/usr/bin/setsid", *subproc,
			"--tag", procname, "--priority", fmt.Sprintf("%d.%d", facility,
				syslog.LOG_CRIT))
		if err != nil {
			return err
		}
	}
	level_val, err := space_log.LevelFromString(*level)
	if err != nil {
		return err
	}
	if level_val != space_log.DefaultLevel {
		space_log.SetLevel(nil, level_val)
	}
	if *filter != "" {
		re, err := regexp.Compile(*filter)
		if err != nil {
			return err
		}
		space_log.SetLevel(re, space_log.LogLevel(math.MinInt32))
	}
	var t *template.Template
	if *format != "" {
		var err error
		t, err = template.New("user").Funcs(funcmap).Parse(*format)
		if err != nil {
			return err
		}
	}
	var textout space_log.TextOutput
	switch strings.ToLower(*output) {
	case "syslog":
		w, err := space_log.NewSyslogOutput(facility, procname)
		if err != nil {
			return err
		}
		if t == nil {
			t = space_log.SyslogTemplate
		}
		textout = w
	case "stdout":
		if t == nil {
			t = space_log.ColorTemplate
		}
		textout = space_log.NewWriterOutput(os.Stdout)
	case "stderr":
		if t == nil {
			t = space_log.ColorTemplate
		}
		textout = space_log.NewWriterOutput(os.Stderr)
	default:
		if t == nil {
			t = space_log.StandardTemplate
		}
		fh, err := os.OpenFile(*output,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		textout = space_log.NewWriterOutput(fh)
	}
	if *buffer > 0 {
		textout = space_log.NewBufferedOutput(textout, *buffer)
	}
	space_log.SetHandler(nil, space_log.NewTextHandler(t, textout))
	log.SetFlags(log.Lshortfile)
	stdlog_level_val, err := space_log.LevelFromString(*stdlog_level)
	if err != nil {
		return err
	}
	log.SetOutput(stdlog.WriterWithoutCaller(stdlog_level_val))
	return nil
}
