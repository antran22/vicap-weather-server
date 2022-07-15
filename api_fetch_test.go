package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIFetcher(t *testing.T) {
	apiFetcher := NewAPIFetcher()
	fetchedResponse := apiFetcher.Fetch(context.Background(), "Da Lat")

	assert.NotNil(t, fetchedResponse)
	assert.NotEmpty(t, fetchedResponse.ResolvedAddress)
	assert.NotEmpty(t, fetchedResponse.Days)
}
