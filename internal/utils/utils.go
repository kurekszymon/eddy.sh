package utils

import (
	"fmt"
	"os"

	"github.com/kurekszymon/eddy.sh/internal/logger"
)

func PromptConfirm(prompt string, error_message string, codes ...int) {
	logger.Prompt(prompt)
	logger.Prompt("Type [Y] or [y] to continue")

	var i string
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		if len(error_message) > 0 {

			logger.Error(error_message)
		}

		if len(codes) > 0 {
			os.Exit(codes[0])
		}
	}
}
