/*
Copyright 2015 All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
