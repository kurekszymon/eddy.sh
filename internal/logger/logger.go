package logger

import (
	"fmt"
)

type LogType string

const (
	LogDebug   LogType = "debug"
	LogWarning LogType = "warning"
	LogError   LogType = "error"
	LogInfo    LogType = "info"
)

func FormatLogType(message string, logType LogType) string {
	const (
		gray   = "\033[90m"
		debug  = "\033[34m"
		yellow = "\033[33m"
		red    = "\033[31m"
		reset  = "\033[0m"
	)

	var color string
	switch logType {
	case LogDebug:
		color = debug
	case LogWarning:
		color = yellow
	case LogError:
		color = red
	case LogInfo:
		color = gray
	default:
		color = reset
	}

	fmt_type := fmt.Sprintf("%s[%s]%s", color, logType, reset)

	fmt_message := fmt.Sprintf("[eddy.sh]%s: %s", fmt_type, message)
	return fmt_message
}

func Info(message string) {
	fmt_message := FormatLogType(message, LogInfo)
	fmt.Println(fmt_message)
}

func Warn(message string) {
	fmt_message := FormatLogType(message, LogWarning)
	fmt.Println(fmt_message)
}

func Error(message string) {
	fmt_message := FormatLogType(message, LogError)
	fmt.Println(fmt_message)
}

func Debug(message string) {
	fmt_message := FormatLogType(message, LogDebug)
	fmt.Println(fmt_message)
}
