package logger

import (
	"io"
	"log"
	"runtime"
	"strings"
)

var debug bool

var debugLogger   *log.Logger
var infoLogger    *log.Logger
var errLogger     *log.Logger

func InitLogger(debugHandle, infoHandle, errorHandle io.Writer, dbg bool) {

	debug = dbg

	debugLogger = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate | log.Ltime)

	infoLogger = log.New(infoHandle,
		"INFO: ",
		log.Ldate | log.Ltime)

	errLogger = log.New(errorHandle,
		"ERROR: ",
		log.Ldate | log.Ltime)
}

func LogDebug(v ...interface{}) {
	if debug {
		debugLogger.Println(v)
	}
}

func LogInfo(v ...interface{}) {
	infoLogger.Println(v)
}

func LogError(v ...interface{}) {
	//if len(v) == 1 && v[0].(error) == nil {
	//	return
	//}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		file = strings.SplitAfter(file, "stock/")[1]
	}
	errLogger.Println(file, ":", line, v)
}


