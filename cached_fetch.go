package main

import (
	"context"
	"log"
)

type CachedAPIFetcher struct {
	apiFetcher *APIFetcher
	store      MetricDataStore
}

func NewCachedAPIFetcher() *CachedAPIFetcher {
	return &CachedAPIFetcher{
		apiFetcher: NewAPIFetcher(),
		store:      NewRedisMetricDataStore(),
	}
}

func (c CachedAPIFetcher) Fetch(ctx context.Context, locationQuery string) *LocalMetrics {
	cachedData := c.store.GetMetrics(ctx, locationQuery)
	if cachedData != nil {
		log.Printf("found cached data for query %s\n", locationQuery)
		return cachedData
	}

	canonicalName, err := c.apiFetcher.ProbeCanonicalName(ctx, locationQuery)

	if err == nil {
		log.Printf("found canonical name for query '%s': '%s'", locationQuery, canonicalName)
		canonicalCachedData := c.store.GetMetrics(ctx, canonicalName)
		log.Printf("found canonical cached data for query: %s; canonical name: %s", locationQuery, canonicalName)

		if canonicalCachedData != nil {
			c.store.BlindStore(ctx, locationQuery, canonicalCachedData)
			return canonicalCachedData
		}
	} else {
		log.Println("unable to fetch canonical name", err)
	}

	metrics := c.apiFetcher.Fetch(ctx, locationQuery)

	if metrics != nil {
		if err = c.store.StoreMetrics(ctx, metrics); err != nil {
			c.store.BlindStore(ctx, locationQuery, metrics)
		}
	}

	return metrics

}
