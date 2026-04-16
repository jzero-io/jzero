package serverless

import (
	"reflect"
	"testing"
)

func TestServerlessDisplayItems(t *testing.T) {
	got := serverlessDisplayItems([]string{
		"plugins/helloworld",
		"go.work",
		"go.work.sum",
		"plugins/plugins.go",
	})

	want := []string{
		"plugins/helloworld",
		"plugins/plugins.go",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("serverlessDisplayItems() = %v, want %v", got, want)
	}
}
