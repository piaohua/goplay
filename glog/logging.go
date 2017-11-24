package glog

import (
	"flag"
	"os"

	logging "go-logging"
)

var log = logging.MustGetLogger("goplay")

var format1 = logging.MustStringFormatter(
	`[%{level:.4s}] [%{time:2006-01-02 15:04:05.000}] [%{shortfile}] %{shortfunc} %{pid} %{message}`,
)

var format2 = logging.MustStringFormatter(
	`%{color}[%{level:.4s}] [%{time:15:04:05.000}] [%{shortfile}] %{pid}%{color:reset} %{message}`,
)

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

//TODO 设置log level
func Init() {
	flag.Parse()
	backend1 := logging.NewLogBackendFile(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend1Formatter := logging.NewBackendFormatter(backend1, format1)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	backend2Formatter := logging.NewBackendFormatter(backend2, format2)
	backend2Leveled := logging.AddModuleLevel(backend2Formatter)
	backend2Leveled.SetLevel(logging.DEBUG, "")

	logging.SetBackend(backend1Leveled, backend2Leveled)
	//logging.SetBackend(backend2Leveled)
}

// Flush flushes all pending log I/O.
func Flush() {
	logging.Flush()
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level.
func Critical(args ...interface{}) {
	log.Critical(args...)
}

// Criticalf logs a message using CRITICAL as log level.
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level.
func Notice(args ...interface{}) {
	log.Notice(args...)
}

// Noticef logs a message using NOTICE as log level.
func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	//TODO
	//if !strings.HasSuffix(format, "\n") {
	//	format += "\n"
	//}
	log.Debugf(format, args...)
}
