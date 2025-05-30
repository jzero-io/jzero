package stringx

import (
	"strings"
)

func FirstUpper(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	return s
}

func FirstLower(s string) string {
	if len(s) > 0 {
		return strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

func ToCamel(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, "/", "-")
	words := strings.Split(s, "-")

	for i := 1; i < len(words); i++ {
		words[i] = FirstUpper(words[i])
	}

	result := strings.Join(words, "")

	return result
}
