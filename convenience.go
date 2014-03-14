// Copyright (C) 2014 Space Monkey, Inc.

package log

import (
	"fmt"
	"io"
)

func (l *Logger) Debug(v ...interface{}) {
	if l.getLevel() <= Debug {
		l.getHandler().Log(l.name, Debug, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.getLevel() <= Debug {
		l.getHandler().Log(l.name, Debug, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) DebugEnabled() bool {
	return l.getLevel() <= Debug
}

func (l *Logger) Info(v ...interface{}) {
	if l.getLevel() <= Info {
		l.getHandler().Log(l.name, Info, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.getLevel() <= Info {
		l.getHandler().Log(l.name, Info, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) InfoEnabled() bool {
	return l.getLevel() <= Info
}

func (l *Logger) Warn(v ...interface{}) {
	if l.getLevel() <= Warning {
		l.getHandler().Log(l.name, Warning, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.getLevel() <= Warning {
		l.getHandler().Log(l.name, Warning, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) WarnEnabled() bool {
	return l.getLevel() <= Warning
}

func (l *Logger) Error(v ...interface{}) {
	if l.getLevel() <= Error {
		l.getHandler().Log(l.name, Error, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.getLevel() <= Error {
		l.getHandler().Log(l.name, Error, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) ErrorEnabled() bool {
	return l.getLevel() <= Error
}

func (l *Logger) Crit(v ...interface{}) {
	if l.getLevel() <= Critical {
		l.getHandler().Log(l.name, Critical, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Critf(format string, v ...interface{}) {
	if l.getLevel() <= Critical {
		l.getHandler().Log(l.name, Critical, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) CritEnabled() bool {
	return l.getLevel() <= Critical
}

func (l *Logger) Log(level LogLevel, v ...interface{}) {
	if l.getLevel() <= level {
		l.getHandler().Log(l.name, level, fmt.Sprint(v...), 1)
	}
}

func (l *Logger) Logf(level LogLevel, format string, v ...interface{}) {
	if l.getLevel() <= level {
		l.getHandler().Log(l.name, level, fmt.Sprintf(format, v...), 1)
	}
}

func (l *Logger) LevelEnabled(level LogLevel) bool {
	return l.getLevel() <= level
}

type writer struct {
	l     *Logger
	level LogLevel
}

func (w *writer) Write(data []byte) (int, error) {
	if w.l.getLevel() <= w.level {
		w.l.getHandler().Log(w.l.name, w.level, string(data), 1)
	}
	return len(data), nil
}

func (l *Logger) Writer(level LogLevel) io.Writer {
	return &writer{l: l, level: level}
}

type writerNoCaller struct {
	l     *Logger
	level LogLevel
}

func (w *writerNoCaller) Write(data []byte) (int, error) {
	if w.l.getLevel() <= w.level {
		w.l.getHandler().Log(w.l.name, w.level, string(data), -1)
	}
	return len(data), nil
}

func (l *Logger) WriterWithoutCaller(level LogLevel) io.Writer {
	return &writerNoCaller{l: l, level: level}
}
