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

func PrintInstallErrors(errors_group ...map[string]error) {
	for _, errors := range errors_group {
		if len(errors) > 0 {
			for toolName, err := range errors {
				message := fmt.Sprintf("Error installing %s: %v\n", toolName, err)
				logger.Error(message)
			}
		}
	}
}
