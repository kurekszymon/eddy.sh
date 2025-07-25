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
	LogPrompt  LogType = "prompt"
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
	case LogPrompt:
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

	formattedType := fmt.Sprintf("%s[%s]%s", color, logType, reset)

	formattedMessage := fmt.Sprintf("[eddy.sh]%s: %s", formattedType, message)
	return formattedMessage
}

func Info(message string) {
	formattedMessage := FormatLogType(message, LogInfo)
	fmt.Println(formattedMessage)
}

func Warn(message string) {
	formattedMessage := FormatLogType(message, LogWarning)
	fmt.Println(formattedMessage)
}

func Error(message string) {
	formattedMessage := FormatLogType(message, LogError)
	fmt.Println(formattedMessage)
}

func Debug(message string) {
	formattedMessage := FormatLogType(message, LogDebug)
	fmt.Println(formattedMessage)
}

func Prompt(message string) {
	formattedMessage := FormatLogType(message, LogPrompt)
	fmt.Println(formattedMessage)
}
