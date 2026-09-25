package utils

import (
	"fmt"
	"os"
	"strings"

	"ghosttools/internal/ui"
)

func RequireInput(prompt string) string {
	ui.Cyan(prompt)
	var s string
	fmt.Scanln(&s)
	return strings.TrimSpace(s)
}

func WriteFile(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		ui.Red(err.Error())
	}
}