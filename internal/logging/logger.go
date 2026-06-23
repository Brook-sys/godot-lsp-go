package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	debug bool
	quiet bool
	out   io.Writer
	file  *os.File
	mu    sync.Mutex
}

func New(debug, quiet bool, logFile string) (*Logger, error) {
	l := &Logger{debug: debug, quiet: quiet, out: os.Stderr}
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		l.file = f
	}
	return l, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Info(format string, args ...any) {
	if l.quiet {
		return
	}
	l.write("info", format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	if !l.debug {
		return
	}
	l.write("debug", format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	if l.quiet {
		return
	}
	l.write("warn", format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.write("error", format, args...)
}

func (l *Logger) write(level, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	line := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format(time.RFC3339), level, fmt.Sprintf(format, args...))
	_, _ = io.WriteString(l.out, line)
	if l.file != nil {
		_, _ = io.WriteString(l.file, line)
	}
}
