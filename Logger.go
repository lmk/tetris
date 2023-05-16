package main

import (
	"io"
	"log"
)

// Trace log
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLogger 로거 초기화
func InitLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	flagTrace := 0
	flagInfo := 0
	flagWarning := 0
	flagError := 0

	if conf.Log.Datetime {
		flagTrace |= log.Ldate | log.Ltime | log.Lshortfile
		flagInfo |= log.Ldate | log.Ltime | log.Lshortfile
		flagWarning |= log.Ldate | log.Ltime | log.Lshortfile
		flagError |= log.Ldate | log.Ltime | log.Lshortfile
	}

	if conf.Log.SrcFile {
		flagTrace |= log.Lshortfile
		flagInfo |= log.Lshortfile
		flagWarning |= log.Lshortfile
		flagError |= log.Lshortfile
	}

	Trace = log.New(traceHandle, "TRACE: ", flagTrace)
	Info = log.New(infoHandle, "INFO: ", flagInfo)
	Warning = log.New(warningHandle, "WARNING: ", flagWarning)
	Error = log.New(errorHandle, "ERROR: ", flagError)

}
