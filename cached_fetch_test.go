package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedAPIFetcher(t *testing.T) {
	fetcher := NewCachedAPIFetcher()
	t.Run("test single query", func(t *testing.T) {
		fetcher.store.Clear()

		firstQuery := fetcher.Fetch(context.Background(), "Da Lat")

		assert.NotNil(t, firstQuery)
	})

	t.Run("test duplicated query in 100ms window", func(t *testing.T) {
		fetcher.store.Clear()

		firstQuery := fetcher.Fetch(context.Background(), "Da Lat")
		assert.NotNil(t, firstQuery)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		secondQuery := fetcher.Fetch(ctx, "Da Lat")

		assert.NotErrorIs(t, ctx.Err(), context.Canceled)
		assert.Equal(t, firstQuery, secondQuery)

	})

	t.Run("test query with similar canonical name in 1s window", func(t *testing.T) {
		fetcher.store.Clear()

		firstQuery := fetcher.Fetch(context.Background(), "Da Lat")
		assert.NotNil(t, firstQuery)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		secondQuery := fetcher.Fetch(ctx, "Đà Lạt")

		assert.NotErrorIs(t, ctx.Err(), context.Canceled)

		assert.Equal(t, firstQuery, secondQuery)
	})

	t.Run("test query with similar canonical name in 1s window", func(t *testing.T) {
		fetcher.store.Clear()

		firstQuery := fetcher.Fetch(context.Background(), "Da Lat")
		assert.NotNil(t, firstQuery)

		secondQuery := fetcher.Fetch(context.Background(), "Hà Nội")

		assert.NotEqual(t, firstQuery, secondQuery)
	})

}
