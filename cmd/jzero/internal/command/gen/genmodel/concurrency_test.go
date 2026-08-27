package genmodel

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

func TestSetModelConcurrencyLimit(t *testing.T) {
	original := config.C.Gen.ModelConcurrency
	t.Cleanup(func() {
		config.C.Gen.ModelConcurrency = original
	})

	config.C.Gen.ModelConcurrency = 2
	var eg errgroup.Group
	setModelConcurrencyLimit(&eg, 5)

	var current atomic.Int32
	var maximum atomic.Int32
	for range 5 {
		eg.Go(func() error {
			running := current.Add(1)
			for {
				observed := maximum.Load()
				if running <= observed || maximum.CompareAndSwap(observed, running) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			current.Add(-1)
			return nil
		})
	}

	require.NoError(t, eg.Wait())
	require.Equal(t, int32(2), maximum.Load())
}

func TestSetModelConcurrencyLimitUnlimited(t *testing.T) {
	original := config.C.Gen.ModelConcurrency
	t.Cleanup(func() {
		config.C.Gen.ModelConcurrency = original
	})

	config.C.Gen.ModelConcurrency = 0
	var eg errgroup.Group
	setModelConcurrencyLimit(&eg, 3)

	var started atomic.Int32
	release := make(chan struct{})
	for range 3 {
		eg.Go(func() error {
			started.Add(1)
			<-release
			return nil
		})
	}

	require.Eventually(t, func() bool {
		return started.Load() == 3
	}, time.Second, time.Millisecond)
	close(release)
	require.NoError(t, eg.Wait())
}
