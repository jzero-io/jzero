package fuzzy

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/microcosm-cc/bluemonday"
)

var (
	once                sync.Once
	EnableXssProtection bool
	BlueMondayPolicy    = bluemonday.StrictPolicy()
)

func Decode(bodyBytes []byte, req any) ([]byte, error) {
	once.Do(func() {
		RegisterFuzzyDecoders()
		RegisterPointerFuzzyDecoders()
	})

	if err := jsoniter.Unmarshal(bodyBytes, &req); err != nil {
		return nil, err
	}

	fuzzyDecodeBytes, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}

	return fuzzyDecodeBytes, nil
}
