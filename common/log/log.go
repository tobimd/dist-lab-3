package log

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	LstdAppendFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	LstdLogFlags    = log.LstdFlags | log.Lmicroseconds
)

var (
	fallbackFileName = "unnamed-log"
)

// Print formatted message to stdout as well as to log file
func Print(f *string, msg string, a ...interface{}) {
	log.Printf(fmt.Sprintf(msg, a...))
	Log(f, msg, a...)
}

// Append formatted message to log file
func Log(f *string, msg string, a ...interface{}) {
	if f == nil {
		f = &fallbackFileName
	}

	file, err := os.OpenFile(*f, LstdAppendFlags, 0600)
	FailOnError(f, err, "Couldn't open file \"%s\"", *f)
	defer file.Close()

	logger := log.New(file, "", LstdLogFlags)
	logger.Println(fmt.Sprintf(msg, a...))
}

// If `err` is not `nil`, then follow with Fatal
func FailOnError(f *string, err error, msg string, a ...interface{}) {
	if err != nil {
		Fatal(f, "%s (error:%v)", fmt.Sprintf(msg, a...), err)
	}
}

// Terminate process printing a formatted message to stdout and
// to log file
func Fatal(f *string, msg string, a ...interface{}) {
	Print(f, "[ FATAL ! ] %s", fmt.Sprintf(msg, a...))
	syscall.Exit(1)
}
