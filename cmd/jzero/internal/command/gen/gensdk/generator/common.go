package generator

import (
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/vars"
)

func getResources(resource vars.ResourceHTTPInterfaceMap) []string {
	var resources []string

	for k := range resource {
		if k != "" { // 排除未分组的API
			resources = append(resources, string(k))
		}
	}

	return resources
}

// getUngroupedAPIs 获取未分组的API
func getUngroupedAPIs(resource vars.ResourceHTTPInterfaceMap) []*vars.HTTPInterface {
	if apis, exists := resource[""]; exists {
		return apis
	}
	return nil
}
