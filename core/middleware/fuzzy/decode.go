package fuzzy

import (
	"sync"

	"github.com/jaronnie/genius"
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

	g, err := genius.NewFromRawJSON(bodyBytes)
	if err != nil {
		return nil, err
	}
	keys := g.GetAllKeys()

	if err := jsoniter.Unmarshal(bodyBytes, &req); err != nil {
		return nil, err
	}

	fuzzyDecodeBytes, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}

	ng, err := genius.NewFromRawJSON(fuzzyDecodeBytes)
	if err != nil {
		return nil, err
	}
	nkeys := ng.GetAllKeys()

	if len(keys) != len(nkeys) {
		for _, key := range keys {
			if err = g.Set(key, ng.Get(key)); err != nil {
				return nil, err
			}
		}
		return g.EncodeToJSON()
	}
	return fuzzyDecodeBytes, nil
}
