package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricDataStore(t *testing.T) {
	storesUnderTest := []struct {
		store MetricDataStore
		name  string
	}{
		{store: NewInMemoryMetricDataStore(), name: "in memory"},
		{store: NewRedisMetricDataStore(), name: "redis"},
	}

	metrics1 := loadSampleAPIResponse(t)

	metrics2 := &LocalMetrics{
		ResolvedAddress: "Hanoi, Vietnam",
		Address:         "Hà Nội",
	}

	ctx := context.Background()

	for _, storeSuite := range storesUnderTest {
		store := storeSuite.store
		t.Run(fmt.Sprintf("testing %s store", storeSuite.name), func(t *testing.T) {
			t.Run("saving and loading", func(t *testing.T) {
				store.Clear()

				assert.NoError(t, store.StoreMetrics(ctx, metrics1))

				assert.Equal(t, metrics1, store.GetMetrics(ctx, metrics1.ResolvedAddress))
				assert.Equal(t, metrics1, store.GetMetrics(ctx, metrics1.Address))
			})

			t.Run("distinct key", func(t *testing.T) {
				store.Clear()

				assert.NoError(t, store.StoreMetrics(ctx, metrics1))
				assert.NoError(t, store.StoreMetrics(ctx, metrics2))

				assert.NotEqual(t, metrics1, store.GetMetrics(ctx, metrics2.ResolvedAddress),
					"distinct key must yield distinct data")

			})

			t.Run("non existing key", func(t *testing.T) {
				store.Clear()
				assert.Nil(t, store.GetMetrics(ctx, "bogus key"))

				assert.NoError(t, store.StoreMetrics(ctx, metrics1))

				assert.Nil(t, store.GetMetrics(ctx, "bogus key"))
			})

			t.Run("aliasing", func(t *testing.T) {
				store.Clear()

				assert.NoError(t, store.StoreMetrics(ctx, metrics1))

				assert.NoError(t, store.Alias(ctx, "new key", metrics1.ResolvedAddress))

				assert.Equal(t, metrics1, store.GetMetrics(ctx, "new key"))
			})
		})

	}
}
