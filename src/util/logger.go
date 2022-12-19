package util

import (
	"fmt"
)

const (
	INFO  = 0
	WARN  = 1
	ERROR = 2
	TRACE = 3
)

type LogLevel int

func internal(prefix string, suffix string, text string, format ...any) {
	fmt.Printf(fmt.Sprintf("\033[34mphantom | im-server\033[0m | %s -> \033[35m%s\033[0m: %s", prefix, suffix, text), format...)
	fmt.Println()
}

func Log(level LogLevel, suffix string, text string, format ...any) {
	var prefix string

	switch level {
	case INFO:
		prefix = "\033[32mINFO\033[0m"
	case WARN:
		prefix = "\033[33mWARN\033[0m"
	case ERROR:
		prefix = "\033[31mERROR\033[0m"
	case TRACE:
		prefix = "\033[36mTRACE\033[0m"
	}

	internal(prefix, suffix, text, format...)
}
