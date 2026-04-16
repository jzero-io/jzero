package progress

import (
	"errors"
	"testing"
)

func TestItemFromError(t *testing.T) {
	err := errors.New("find route api files: parse api file: desc/api/user.api, err: user.api 12:2 syntax error")
	if got, want := ItemFromError(err), "desc/api/user.api"; got != want {
		t.Fatalf("ItemFromError() = %q, want %q", got, want)
	}
}
