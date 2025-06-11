package generator

import (
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/vars"
)

func getResources(resource vars.ResourceHTTPInterfaceMap) []string {
	var resources []string

	for k := range resource {
		resources = append(resources, string(k))
	}

	return resources
}
