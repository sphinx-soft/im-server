package logging

import "fmt"

func internal(prefix string, suffix string, text string, format ...any) {
	fmt.Printf(fmt.Sprintf("\033[34mchimera | im-server\033[0m | %s -> \033[35m%s\033[0m: %s", prefix, suffix, text), format...)
	fmt.Println()
}

func Info(suffix string, text string, format ...any) {
	internal("\033[32m  INFO\033[0m", suffix, text, format...)
}

func Warn(suffix string, text string, format ...any) {
	internal("\033[33m  WARN\033[0m", suffix, text, format...)
}

func Error(suffix string, text string, format ...any) {
	internal("\033[31m ERROR\033[0m", suffix, text, format...)
}

func Fatal(suffix string, text string, format ...any) {
	internal("\033[38;5;89m FATAL\033[0m", suffix, text, format...)
}

func Trace(suffix string, text string, format ...any) {
	internal("\033[36m TRACE\033[0m", suffix, text, format...)
}

func Debug(suffix string, text string, format ...any) {
	internal("\033[38;5;99m DEBUG\033[0m", suffix, text, format...)
}

func System(suffix string, text string, format ...any) {
	internal("\033[38;5;172mSYSTEM\033[0m", suffix, text, format...)
}
