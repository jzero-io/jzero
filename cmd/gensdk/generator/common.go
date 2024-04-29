package generator

import "github.com/jzero-io/jzero/cmd/gensdk/vars"

func getScopes(rhis vars.ScopeResourceHTTPInterfaceMap) []string {
	var scopes []string
	for k := range rhis {
		scopes = append(scopes, string(k))
	}

	return scopes
}

func getScopeResources(resource vars.ResourceHTTPInterfaceMap) []string {
	var resources []string

	for k := range resource {
		resources = append(resources, string(k))
	}

	return resources
}
