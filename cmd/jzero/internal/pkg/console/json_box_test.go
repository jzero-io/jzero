package console

import (
	"regexp"
	"strings"
	"testing"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestJSONBoxLines(t *testing.T) {
	lines, err := JSONBoxLines("Parse Config", []byte(`{"debug":true,"name":"demo","items":[1,"x"]}`))
	if err != nil {
		t.Fatalf("JSONBoxLines() error = %v", err)
	}

	if len(lines) < 5 {
		t.Fatalf("JSONBoxLines() returned too few lines: %v", lines)
	}

	if got, want := stripANSI(lines[0]), "┌─  Parse Config"; got != want {
		t.Fatalf("JSONBoxLines()[0] = %q, want %q", got, want)
	}

	if got := stripANSI(lines[2]); got != `│      "debug": true,` {
		t.Fatalf("JSONBoxLines()[2] = %q", got)
	}

	if got := stripANSI(lines[len(lines)-2]); got != "└─ ✓ Complete" {
		t.Fatalf("JSONBoxLines() footer = %q", got)
	}

	if !strings.Contains(lines[2], "\x1b[") {
		t.Fatalf("JSONBoxLines()[2] should contain ANSI styling, got %q", lines[2])
	}
}

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}
