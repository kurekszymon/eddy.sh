package utils

import (
	"fmt"
	"os"
)

func PromptConfirm(prompt string, error_message string, codes ...int) {
	fmt.Print(prompt)

	var i string
	fmt.Scan(&i)

	if i != "Y" && i != "y" {
		fmt.Println(error_message)

		if len(codes) > 0 {
			os.Exit(codes[0])
		}
	}
}
