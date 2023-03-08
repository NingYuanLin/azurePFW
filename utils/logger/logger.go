package logger

import (
	"log"
	"os"
)

const (
	defaultFlag     = log.Ldate | log.Ltime | log.Lshortfile
	defaultPreDebug = "[DEBUG]"
	defaultPreInfo  = "[INFO]"
	defaultPreWarn  = "[WARN]"
	defaultPreError = "[ERROR]"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	debugLogger = log.New(os.Stdout, defaultPreDebug, defaultFlag)
	infoLogger = log.New(os.Stdout, defaultPreInfo, defaultFlag)
	warnLogger = log.New(os.Stdout, defaultPreWarn, defaultFlag)
	errorLogger = log.New(os.Stdout, defaultPreError, defaultFlag)
}

func Debug(v ...any) {
	debugLogger.Print(v...)
}

func Info(v ...any) {
	infoLogger.Print(v...)
}

func Warn(v ...any) {
	warnLogger.Print(v...)
}

func Error(v ...any) {
	errorLogger.Print(v...)
}

func Debugf(format string, v ...any) {
	debugLogger.Printf(format, v...)
}

func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Warnf(format string, v ...any) {
	warnLogger.Printf(format, v...)
}

func Errorf(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

func SetDebugOutPutPath(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	debugLogger.SetOutput(file)
}

func SetInfoOutPutPath(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	infoLogger.SetOutput(file)
}

func SetWarnOutPutPath(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	warnLogger.SetOutput(file)
}

func SetErrorOutPutPath(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	errorLogger.SetOutput(file)
}

func SetDebugPrefix(prefix string) {
	debugLogger.SetPrefix(prefix)
}

func SetInfoPrefix(prefix string) {
	infoLogger.SetPrefix(prefix)
}

func SetWarnPrefix(prefix string) {
	warnLogger.SetPrefix(prefix)
}

func SetErrorPrefix(prefix string) {
	errorLogger.SetPrefix(prefix)
}

func SetDebugFlag(flag int) {
	debugLogger.SetFlags(flag)
}

func SetInfoFlag(flag int) {
	infoLogger.SetFlags(flag)
}

func SetWarnFlag(flag int) {
	warnLogger.SetFlags(flag)
}

func SetErrorFlag(flag int) {
	errorLogger.SetFlags(flag)
}
