package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	ERROR Level = iota
	INFO
	DEBUG
)

type Logger struct {
	verbose bool
	out     io.Writer
}

func New(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		out:     os.Stdout,
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.out = w
}

func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	if l.verbose {
		l.log(DEBUG, format, args...)
	}
}

func (l *Logger) Verbose() bool {
	return l.verbose
}

func (l *Logger) log(level Level, format string, args ...any) {
	var prefix string

	switch level {
	case ERROR:
		prefix = "[ERROR]"
	case INFO:
		prefix = "[INFO]"
	case DEBUG:
		prefix = "[DEBUG]"
	}
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "%s %s %s\n", timestamp, prefix, message)
}
