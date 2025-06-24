package utils

import (
	"fmt"
	"os"

	"github.com/kurekszymon/eddy.sh/internal/types"
)

func FormatLogType(message string, logType types.LogType) string {
	const (
		gray   = "\033[90m"
		debug  = "\033[34m"
		yellow = "\033[33m"
		red    = "\033[31m"
		reset  = "\033[0m"
	)

	var color string
	switch logType {
	case types.LogDebug:
		color = debug
	case types.LogWarning:
		color = yellow
	case types.LogError:
		color = red
	case types.LogInfo:
		color = gray
	default:
		color = reset
	}

	fmt_type := fmt.Sprintf("%s[%s]%s", color, logType, reset)

	fmt_message := fmt.Sprintf("[eddy.sh]%s: %s", fmt_type, message)
	return fmt_message
}

func Log(message string, logType types.LogType) {
	fmt_message := FormatLogType(message, logType)
	fmt.Println(fmt_message)
}

func PromptConfirm(prompt string, error_message string, codes ...int) {
	fmt.Print(prompt)

	var i string
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		Log(error_message, types.LogError)

		if len(codes) > 0 {
			os.Exit(codes[0])
		}
	}
}
