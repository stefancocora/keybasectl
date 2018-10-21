// Package log sets up logging
package log

import (
	"io"
	"log"
	"os"
)

// logging vars and settings
var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
)

// LoggingInit function initialises the writers for each type of logging making it easy to switch between os.Stdout/os.Stderr/ioutil.Discard/file ...
// IN  (os.Stdout, "short/long", os.Stderr, "short/long", os.Stderr, "short/long")
// OUT no return
// but sets up the global package logging parameters
func LoggingInit(
	infoHandle io.Writer,
	infoFilenameLength string,
	errorHandle io.Writer,
	errorFilenameLength string,
	debugHandle io.Writer,
	debugFilenameLength string,
) {
	// should also test for allowed io.Writer destinations
	// should only allow: os.Stdout/os.Stderr/ioutil.Discard/filename
	if infoFilenameLength == "short" {
		InfoLog = log.New(infoHandle, "level=info  ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	} else if infoFilenameLength == "long" {
		InfoLog = log.New(infoHandle, "level=info  ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	} else {
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("don't know how to set this log filename length, allowed values are only \"long\" or \"short\" and you've passed  %v\n", infoFilenameLength)
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("stopping execution \n")
		os.Exit(1)
	}

	if errorFilenameLength == "short" {
		ErrorLog = log.New(errorHandle, "level=error ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	} else if errorFilenameLength == "long" {
		ErrorLog = log.New(errorHandle, "level=error ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	} else {
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("don't know how to set this log filename length, allowed values are only \"long\" or \"short\" and you've passed  %v\n", errorFilenameLength)
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("stopping execution \n")
		os.Exit(1)
	}

	if debugFilenameLength == "short" {
		DebugLog = log.New(debugHandle, "level=debug ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	} else if debugFilenameLength == "long" {
		DebugLog = log.New(debugHandle, "level=debug ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	} else {
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("don't know how to set this log filename length, allowed values are only \"long\" or \"short\" and you've passed  %v\n", errorFilenameLength)
		log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile).Printf("stopping execution \n")
		os.Exit(1)
	}
}
