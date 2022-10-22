package util

import (
	"fmt"
)

func Log(prefix string, text string, format ...any) {
	fmt.Printf("[\033[35m"+prefix+"\033[0m] "+text, format...)
	fmt.Println()
}

func Error(text string, format ...any) {
	Log("Error", text, format...)
}

func Debug(text string, format ...any) {
	Log("Debug", text, format...)
}
