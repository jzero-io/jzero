package progress

import (
	"path/filepath"
	"regexp"
	"strings"
)

var errorItemPattern = regexp.MustCompile(`([[:alnum:]_./-]+\.(?:api|sql|proto|go|json))`)

// ItemFromError extracts the first file-like item from an error for progress display.
func ItemFromError(err error) string {
	if err == nil {
		return ""
	}

	match := errorItemPattern.FindString(err.Error())
	if match == "" {
		return ""
	}

	return filepath.ToSlash(strings.TrimRight(match, ",:"))
}
