package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	LstdAppendFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	LstdLogFlags          = log.LstdFlags | log.Lmicroseconds
	outputDir             = ".logs"
	formattedLogOutputDir = outputDir + "/%s"
)

var (
	fallbackFileName = "unnamed.log"
)

// Print formatted message to stdout as well as to log file
func Print(f *string, msg string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf(msg, a...) + "\n")
	Log(f, msg, a...)
}

// Append formatted message to log file
func Log(f *string, msg string, a ...interface{}) {
	// Format file name to have prefix
	if f == nil || *f == "" {
		f = &fallbackFileName
	}
	prefixedFilename := fmt.Sprintf(formattedLogOutputDir, *f)

	// Check if output log dir exists, if not, then create
	if !FileExists(outputDir) {
		err := os.Mkdir(outputDir, 0755)
		FailOnError(f, err, "Couldn't create directory \"%s\" for logs", outputDir)
	}

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
		Print(f, "[ FATAL ! ] %s", fmt.Sprintf(msg, a...))
		Print(f, "[ ERROR ! ] %v\n", err)
		syscall.Exit(1)
	}
}

// Terminate process printing a formatted message to stdout and
// to log file
func Fatal(f *string, msg string, a ...interface{}) {
	Print(f, "[ FATAL ! ] %s\n", fmt.Sprintf(msg, a...))
	syscall.Exit(1)
}

// Return true if file exists, false otherwise
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}
