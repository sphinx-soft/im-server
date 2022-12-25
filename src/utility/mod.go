package utility

import (
	"bytes"
	"strings"
)

func GetBuild() string {
	// Specifier (Major.Minor.Push.Hotfix)
	return "Next Beta 1 (2.0.0.0)"
}

func SanitizeString(str string) string {
	return strings.Replace(string(bytes.Trim([]byte(str), "\x00")), "\r\n", "", -1)
}
