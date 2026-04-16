package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	jsonKeyValuePattern = regexp.MustCompile(`^(\s*)("(?:\\.|[^"])*")(:\s*)(.*)$`)
	jsonStringPattern   = regexp.MustCompile(`^(\s*)("(?:\\.|[^"])*")([,\]}]?)$`)
	jsonNumberPattern   = regexp.MustCompile(`^(\s*)(-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?)(,?)$`)
	jsonLiteralPattern  = regexp.MustCompile(`^(\s*)(true|false|null)(,?)$`)
)

// BoxCodeItem creates a box line for code-like content.
func BoxCodeItem(item string) string {
	return fmt.Sprintf("│    %s", item)
}

// JSONBoxLines renders JSON content as box lines with light syntax highlighting.
func JSONBoxLines(title string, payload []byte) ([]string, error) {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, payload, "", "  "); err != nil {
		return nil, err
	}

	lines := []string{BoxHeader("", title)}
	for _, line := range strings.Split(pretty.String(), "\n") {
		lines = append(lines, BoxCodeItem(highlightJSONLine(line)))
	}
	lines = append(lines, BoxSuccessFooter(), "")

	return lines, nil
}

// DisplayJSONBox prints a formatted JSON box to stdout.
func DisplayJSONBox(title string, payload []byte) error {
	lines, err := JSONBoxLines(title, payload)
	if err != nil {
		return err
	}

	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

func highlightJSONLine(line string) string {
	if match := jsonKeyValuePattern.FindStringSubmatch(line); match != nil {
		return match[1] +
			Bold(Yellow(match[2])) +
			DimCyan(match[3]) +
			highlightJSONValue(match[4])
	}

	return highlightJSONValue(line)
}

func highlightJSONValue(value string) string {
	if match := jsonStringPattern.FindStringSubmatch(value); match != nil {
		return match[1] + Green(match[2]) + DimCyan(match[3])
	}

	if match := jsonNumberPattern.FindStringSubmatch(value); match != nil {
		return match[1] + Cyan(match[2]) + DimCyan(match[3])
	}

	if match := jsonLiteralPattern.FindStringSubmatch(value); match != nil {
		return match[1] + Yellow(match[2]) + DimCyan(match[3])
	}

	trimmed := strings.TrimSpace(value)
	switch trimmed {
	case "{", "}", "[", "]", "},", "],":
		return strings.TrimSuffix(value, trimmed) + Cyan(trimmed)
	default:
		return value
	}
}
