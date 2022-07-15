package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v9"
)

type RedisMetricDataStore struct {
	client *redis.Client
}

func NewRedisMetricDataStore() *RedisMetricDataStore {
	address := LoadRequiredEnv("REDIS_HOST")
	return &RedisMetricDataStore{
		client: redis.NewClient(
			&redis.Options{
				Addr:     address,
				DB:       0,
				Password: "",
			},
		),
	}
}

func (r *RedisMetricDataStore) StoreMetrics(ctx context.Context, metrics *LocalMetrics) error {
	metricsByte, err := json.Marshal(*metrics)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, metrics.ResolvedAddress, metricsByte, 0).Err()
	if err != nil {
		return err
	}

	err = r.Alias(ctx, metrics.Address, metrics.ResolvedAddress)

	return err
}

func (r *RedisMetricDataStore) BlindStore(ctx context.Context, key string, metrics *LocalMetrics) {
	metricsByte, err := json.Marshal(*metrics)
	if err == nil {
		r.client.Set(ctx, metrics.ResolvedAddress, metricsByte, 0)
	}
}

func (r *RedisMetricDataStore) GetMetrics(ctx context.Context, location string) *LocalMetrics {
	key := location
	for {
		storedBytes, err := r.client.Get(ctx, key).Bytes()
		if err != nil || len(storedBytes) == 0 {
			return nil
		}

		if bytes.HasPrefix(storedBytes, []byte("ALIAS")) {
			key = strings.TrimPrefix(string(storedBytes), "ALIAS")
		} else {
			result := LocalMetrics{}
			err := json.Unmarshal(storedBytes, &result)
			if err != nil {
				return nil
			}
			return &result
		}
	}
}

func (r *RedisMetricDataStore) Alias(ctx context.Context, newKey string, targetKey string) error {
	return r.client.Set(ctx, newKey, "ALIAS"+targetKey, 0).Err()
}

func (r *RedisMetricDataStore) Clear() {
	r.client.FlushDB(context.Background())
}
