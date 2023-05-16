package main

import (
	"io"
	"log"
)

// Trace log
var (
	Debug   *log.Logger
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLogger 로거 초기화
func InitLogger(
	debugHandle io.Writer,
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	flagDebug := 0
	flagTrace := 0
	flagInfo := 0
	flagWarning := 0
	flagError := 0

	if conf.Log.Datetime {
		flagDebug |= log.Ldate | log.Ltime | log.Lshortfile
		flagTrace |= log.Ldate | log.Ltime | log.Lshortfile
		flagInfo |= log.Ldate | log.Ltime | log.Lshortfile
		flagWarning |= log.Ldate | log.Ltime | log.Lshortfile
		flagError |= log.Ldate | log.Ltime | log.Lshortfile
	}

	if conf.Log.SrcFile {
		flagDebug |= log.Lshortfile
		flagTrace |= log.Lshortfile
		flagInfo |= log.Lshortfile
		flagWarning |= log.Lshortfile
		flagError |= log.Lshortfile
	}

	Debug = log.New(debugHandle, "DBG: ", flagDebug)
	Trace = log.New(traceHandle, "TRC: ", flagTrace)
	Info = log.New(infoHandle, "INF: ", flagInfo)
	Warning = log.New(warningHandle, "WAR: ", flagWarning)
	Error = log.New(errorHandle, "ERR: ", flagError)

	if !conf.Log.Debug {
		Debug.SetOutput(io.Discard)
	}

	if !conf.Log.Trace {
		Trace.SetOutput(io.Discard)
	}

	if !conf.Log.Info {
		Info.SetOutput(io.Discard)
	}

	if !conf.Log.Warning {
		Warning.SetOutput(io.Discard)
	}

	if !conf.Log.Error {
		Error.SetOutput(io.Discard)
	}
}
