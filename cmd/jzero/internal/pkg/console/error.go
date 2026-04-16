package console

import (
	"errors"
	"strings"
)

// RenderedError marks an error as already rendered in the console UI.
type RenderedError struct {
	Err error
}

func (e *RenderedError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *RenderedError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// MarkRenderedError wraps an error so the root command won't print it again.
func MarkRenderedError(err error) error {
	if err == nil {
		return nil
	}

	var rendered *RenderedError
	if errors.As(err, &rendered) {
		return err
	}

	return &RenderedError{Err: err}
}

// IsRenderedError reports whether the error has already been rendered.
func IsRenderedError(err error) bool {
	var rendered *RenderedError
	return errors.As(err, &rendered)
}

// NormalizeErrorLines converts free-form error output into unique display lines.
func NormalizeErrorLines(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")

	lines := strings.Split(text, "\n")
	seen := make(map[string]struct{}, len(lines))
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "exit status 1" {
			continue
		}
		line = strings.TrimPrefix(line, "Error: ")
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		result = append(result, line)
	}

	return result
}
