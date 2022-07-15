package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func loadSampleAPIResponse(t testing.TB) *LocalMetrics {
	t.Helper()

	sampleMetrics := LocalMetrics{}
	fileBytes, err := os.ReadFile("sample_api_response.json")
	if err != nil {
		t.Fatal(err)
		return nil
	}

	if err := json.Unmarshal(fileBytes, &sampleMetrics); err != nil {
		return nil
	}

	return &sampleMetrics
}

func LoadRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required env variable '%s'", key))
	}

	return value
}
