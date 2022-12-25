package logging

import "fmt"

func internal(prefix string, suffix string, text string, format ...any) {
	fmt.Printf(fmt.Sprintf("\033[34mphantom | im-server\033[0m | %s -> \033[35m%s\033[0m: %s", prefix, suffix, text), format...)
	fmt.Println()
}

func Info(suffix string, text string, format ...any) {
	internal("\033[32mINFO \033[0m", suffix, text, format...)
}

func Warn(suffix string, text string, format ...any) {
	internal("\033[33mWARN \033[0m", suffix, text, format...)
}

func Error(suffix string, text string, format ...any) {
	internal("\033[31mERROR\033[0m", suffix, text, format...)
}

func Fatal(suffix string, text string, format ...any) {
	internal("\033[38;5;89mFATAL\033[0m", suffix, text, format...)
}

func Trace(suffix string, text string, format ...any) {
	internal("\033[36mTRACE\033[0m", suffix, text, format...)
}

func Debug(suffix string, text string, format ...any) {
	internal("\033[38;5;99mDEBUG\033[0m", suffix, text, format...)
}
