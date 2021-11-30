package log

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	LstdAppendFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	LstdLogFlags          = log.LstdFlags | log.Lmicroseconds
	formattedLogOutputDir = "output_logs/%s"
)

var (
	fallbackFileName = "unnamed.log"
)

// Print formatted message to stdout as well as to log file
func Print(f *string, msg string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf(msg, a...))
	Log(f, msg, a...)
}

// Append formatted message to log file
func Log(f *string, msg string, a ...interface{}) {
	if f == nil {
		f = &fallbackFileName
	}
	prefixedFilename := fmt.Sprintf(formattedLogOutputDir, *f)

	file, err := os.OpenFile(prefixedFilename, LstdAppendFlags, 0600)
	FailOnError(&fallbackFileName, err, "Couldn't open file \"%s\"", prefixedFilename)
	defer file.Close()

	logger := log.New(file, "", LstdLogFlags)
	logger.Println(fmt.Sprintf(msg, a...))
	logger = nil
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
	Print(f, "[ FATAL ! ] %s\n", fmt.Sprintf(msg, a...))
	syscall.Exit(1)
}
