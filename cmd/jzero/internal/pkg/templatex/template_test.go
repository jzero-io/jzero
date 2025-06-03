package templatex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplate(t *testing.T) {
	template, err := ParseTemplate(map[string]any{
		"Style": "go_zero",
	}, []byte(`{{FormatStyle .Style "service_context.go.tpl"}}`))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "service_context.go.tpl", string(template))
}
