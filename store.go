package main

import (
	"context"
)

type MetricDataStore interface {
	StoreMetrics(ctx context.Context, metrics *LocalMetrics) error
	BlindStore(ctx context.Context, key string, metrics *LocalMetrics)
	GetMetrics(ctx context.Context, key string) *LocalMetrics
	Alias(ctx context.Context, newKey string, targetKey string) error
	Clear()
}

type InMemoryMetricDataStore struct {
	mapping map[string]any
}

func (i *InMemoryMetricDataStore) BlindStore(ctx context.Context, key string, metrics *LocalMetrics) {
	i.mapping[key] = metrics
}

func NewInMemoryMetricDataStore() *InMemoryMetricDataStore {
	return &InMemoryMetricDataStore{
		mapping: make(map[string]any),
	}
}

func (i *InMemoryMetricDataStore) StoreMetrics(ctx context.Context, metrics *LocalMetrics) error {
	i.mapping[metrics.ResolvedAddress] = metrics
	i.mapping[metrics.Address] = metrics.ResolvedAddress
	return nil
}

func (i *InMemoryMetricDataStore) GetMetrics(ctx context.Context, location string) *LocalMetrics {
	key := location
	for {
		if value, found := i.mapping[key]; found {
			switch newValue := value.(type) {
			case *LocalMetrics:
				return newValue
			case string:
				key = newValue
			default:
				return nil
			}

		} else {
			return nil
		}
	}
}

func (i *InMemoryMetricDataStore) Clear() {
	i.mapping = make(map[string]any)
}

func (i *InMemoryMetricDataStore) Alias(ctx context.Context, newKey string, targetKey string) error {
	i.mapping[newKey] = targetKey
	return nil
}
