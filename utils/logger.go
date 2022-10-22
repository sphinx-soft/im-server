package utils

import (
	"fmt"
)

func Log(prefix string, text string, format ...any) {
	fmt.Printf("["+prefix+"] "+text, format...)
	fmt.Println()
}
func Error(text string, format ...any) {
	Log("Error", text, format)
}
