package log

import (
	"fmt"
	"log"
)

// Ideas:
// 1. Log everything to a rotating in memory buffer, and dump last n recorded logs on dump
// 2. Dynamically change log levels at will
// 3. Compilation flags for release vs debug
// 4. https://cs.opensource.google/go/go/+/refs/tags/go1.18.2:src/log/log.go;drc=0a1a092c4b56a1d4033372fbd07924dad8cbb50b;l=175

const (
	LevelTrace = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// From: https://pkg.go.dev/log#pkg-constants
const (
	Ldate         = log.Ldate         // the date in the local time zone: 2009/01/23
	Ltime         = log.Ltime         // the time in the local time zone: 01:23:23
	Lmicroseconds = log.Lmicroseconds // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile     = log.Llongfile     // full file name and line number: /a/b/c/d.go:23
	Lshortfile    = log.Lshortfile    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC          = log.LUTC          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix    = log.Lmsgprefix    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime     // initial values for the standard logger
)

// TODO thread safety
var state struct {
	level int
}
func SetLevel(level int) {
	state.level = level
}

func SetFlags(flag int) {
	log.SetFlags(flag)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func Println(v ...any) {
	log.Println(v...)
}

func Trace(format string, v ...any) {
	if state.level > LevelTrace { return }
	log.Printf("[TRACE] " + format, v...)
}

func Info(format string, v ...any) {
	if state.level > LevelInfo { return }
	log.Printf("[INFO]  " + format, v...)
}

func Warn(format string, v ...any) {
	if state.level > LevelWarn { return }
	log.Printf("[WARN]  " + format, v...)
}

func Error(format string, v ...any) {
	if state.level >= LevelError { return }
	log.Printf("[ERROR] " + format, v...)
}

func Fatal(v ...any) {
	log.Fatal(fmt.Sprintf("[FATAL]  ", v...))
}
