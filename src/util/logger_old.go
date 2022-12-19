package util

import (
	"fmt"
)

var printDbg bool = true

func LogOld(prefix string, text string, format ...any) {
	fmt.Printf(fmt.Sprintf("[\033[35m%s\033[0m] %s", prefix, text), format...)
	fmt.Println()
}

func Error(prefix string, text string, format ...any) {
	LogOld("Error", fmt.Sprintf("[\033[31m%s\033[0m] %s", prefix, text), format...)
}

func Debug(prefix string, text string, format ...any) {
	if printDbg {
		LogOld("Debug", fmt.Sprintf("[\033[36m%s\033[0m] %s", prefix, text), format...)
	}
}
