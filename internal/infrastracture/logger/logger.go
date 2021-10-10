package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
)

type logLevel int

const (
	logPanic logLevel = iota
	logError
	logWarning
	logMessage
)

var (
	colorPanic = color.New(color.FgHiWhite).Add(color.BgRed)
	colorError = color.New(color.FgHiRed)
	colorWarning = color.New(color.FgHiYellow)
	colorMessage = color.New(color.FgHiWhite)

	colorCook = color.New(color.FgHiGreen)
	colorSuper = color.New(color.FgHiMagenta)

	prefixPanic = "[PNC] "
	prefixError = "[ERR] "
	prefixWarning = "[WRN] "
	prefixMessage = "[MSG] "

	prefixCook = "[COOK-%d] "
	prefixSuper = "[SPR] "

	LogChannel = make(chan string, 50)
)

func logCustomF(clr *color.Color, prefix string, postLog func(), level logLevel, format string, items ...interface{}){
	if level <= logLevel(viper.GetInt("log_level")) {
		var data string
		if len(items) == 0 {
			data = format
		} else {
			data = fmt.Sprintf(format, items...)
		}

		logData := prefix + data
		LogChannel <- logData

		if postLog != nil {
			postLog()
		}
	}
}

func LogPanicF(format string, items ...interface{}){
	logCustomF(colorPanic, prefixPanic, func(){ panic(format) }, logPanic, format, items...)
}

func LogErrorF(format string, items ...interface{}){
	logCustomF(colorError, prefixError, func(){os.Exit(1)}, logError, format, items...)
}

func LogWarningF(format string, items ...interface{}) {
	logCustomF(colorWarning, prefixWarning, nil, logWarning, format, items...)
}

func LogMessageF(format string, items ...interface{}) {
	logCustomF(colorMessage, prefixMessage, nil, logMessage, format, items...)
}

func LogCookF(id int, format string, items ...interface{}){
	logCustomF(colorCook, fmt.Sprintf(prefixCook, id), nil, logMessage, format, items...)
}

func LogSuperF(format string, items ...interface{}){
	logCustomF(colorSuper, prefixSuper, nil, logMessage, format, items...)
}

func LogPanic(format string){
	LogPanicF(format, nil...)
}

func LogError(format string){
	LogErrorF(format, nil...)
}

func LogWarning(format string) {
	LogWarningF(format, nil...)
}

func LogMessage(format string) {
	LogMessageF(format, nil...)
}

func LogCook(id int, format string) {
	LogCookF(id, format, nil...)
}

func LogSuper(format string) {
	LogSuperF(format, nil...)
}