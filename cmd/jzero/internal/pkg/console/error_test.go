package console

import (
	"errors"
	"testing"
)

func TestMarkRenderedError(t *testing.T) {
	err := errors.New("boom")
	rendered := MarkRenderedError(err)

	if !IsRenderedError(rendered) {
		t.Fatalf("expected rendered error marker")
	}

	if rendered.Error() != err.Error() {
		t.Fatalf("rendered error = %q, want %q", rendered.Error(), err.Error())
	}
}

func TestNormalizeErrorLines(t *testing.T) {
	input := "Error: parse api file: desc/api/user.api\nError: parse api file: desc/api/user.api\n\nexit status 1\nError: find route api files: parse api file: desc/api/user.api"
	got := NormalizeErrorLines(input)

	want := []string{
		"parse api file: desc/api/user.api",
		"find route api files: parse api file: desc/api/user.api",
	}

	if len(got) != len(want) {
		t.Fatalf("NormalizeErrorLines() len = %d, want %d; got=%v", len(got), len(want), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("NormalizeErrorLines()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
