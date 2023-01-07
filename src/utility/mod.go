package utility

import (
	"math/rand"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func GetBuild() string {
	// Specifier (Major.Minor.Push.Hotfix)
	return "Next Beta 1 (2.0.0.0)"
}

// SanitizeString is no longer needed

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomNumber(n int) int {
	return rand.Intn(n)
}

func ConvertToUTF16LE(input string) []byte {
	bytes, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder(), []byte(input))
	if err != nil {
		return nil
	}
	return bytes
}
