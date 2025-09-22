package cache

import "strings"

const keySeparator = ","

func formatKeys(keys []string) string {
	return strings.Join(keys, keySeparator)
}
