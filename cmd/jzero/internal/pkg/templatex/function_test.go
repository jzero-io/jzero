package templatex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStyle(t *testing.T) {
	template, err := ParseTemplate("test", map[string]any{
		"Style": "go_zero",
	}, []byte(`{{FormatStyle .Style "service_context.go.tpl"}}`))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "service_context.go.tpl", string(template))
}

func TestVersionCompare(t *testing.T) {
	template, err := ParseTemplate("test", map[string]any{
		"GoVersion": "1.24",
	}, []byte(`{{if (VersionCompare .GoVersion ">=" "1.24")}}true{{end}}`))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "true", string(template))
}
