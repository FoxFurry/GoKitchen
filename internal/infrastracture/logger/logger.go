package logger

import (
	"github.com/spf13/viper"
	"log"
)

type logLevel int

const (
	logPanic logLevel = iota
	logError
	logWarning
	logMessage
)

func logNoItems(level logLevel, format string){
	if level == logPanic {
		log.Panic(format)
	}else if level <= logLevel(viper.GetInt("log_level")) {
		log.Print(format)
	}
}

func logCustom(level logLevel, format string, items []interface{}){
	if items == nil {
		logNoItems(level, format)
	}else {
		if level == logPanic {
			log.Panicf(format, items)
		}else if level <= logLevel(viper.GetInt("log_level")) {
			log.Printf(format, items)
		}
	}
}

func LogPanicF(format string, items ...interface{}){
	logCustom(logPanic, format, items)
}

func LogErrorF(format string, items ...interface{}){
	logCustom(logError, format, items)
}

func LogWarningF(format string, items ...interface{}) {
	logCustom(logWarning, format, items)
}

func LogMessageF(format string, items ...interface{}) {
	logCustom(logMessage, format, items)
}

func LogPanic(format string){
	logCustom(logPanic, format, nil)
}

func LogError(format string){
	logCustom(logError, format, nil)
}

func LogWarning(format string) {
	logCustom(logWarning, format, nil)
}

func LogMessage(format string) {
	logCustom(logMessage, format, nil)
}