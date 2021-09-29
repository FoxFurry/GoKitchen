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
)


func logCustomf(clr *color.Color, prefix string, postLog func(), level logLevel, format string, items ...interface{}){
	if level <= logLevel(viper.GetInt("log_level")) {
		var data string
		if len(items) == 0 {
			data = format
		} else {
			data = fmt.Sprintf(format, items...)
		}

		clr.Print(prefix + data + "\n")
		if postLog != nil {
			postLog()
		}
	}
}

func LogPanicf(format string, items ...interface{}){
	logCustomf(colorPanic, prefixPanic, func(){ panic(format) }, logPanic, format, items...)
}

func LogErrorF(format string, items ...interface{}){
	logCustomf(colorError, prefixError, func(){os.Exit(1)}, logError, format, items...)
}

func LogWarningF(format string, items ...interface{}) {
	logCustomf(colorWarning, prefixWarning, nil, logWarning, format, items...)
}

func LogMessageF(format string, items ...interface{}) {
	logCustomf(colorMessage, prefixMessage, nil, logMessage, format, items...)
}

func LogCookF(id int, format string, items ...interface{}){
	logCustomf(colorCook, fmt.Sprintf(prefixCook, id), nil, logMessage, format, items...)
}

func LogSuperF(format string, items ...interface{}){
	logCustomf(colorSuper, prefixSuper, nil, logMessage, format, items...)
}

func LogPanic(format string){
	LogPanicf(format, nil...)
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