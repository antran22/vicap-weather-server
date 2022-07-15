package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

type APIFetcher struct {
	key       string
	startDate string
	endDate   string
}

func NewAPIFetcher() *APIFetcher {
	key := LoadRequiredEnv("VISUAL_CROSSING_API_KEY")
	env := os.Getenv("GO_ENV")
	if env == "production" {
		log.Println("Creating API Fetcher in production mode")
		return &APIFetcher{key: key, startDate: "2021-01-01", endDate: "2021-12-31"}
	}

	log.Println("Creating API Fetcher in development mode")
	return &APIFetcher{key: key, startDate: "2021-01-01", endDate: "2021-01-01"}
}

func (A *APIFetcher) Fetch(ctx context.Context, locationQuery string) *LocalMetrics {
	log.Printf("api fetching data for location %s\n", locationQuery)
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s/%s?key=%s&elements=datetime,hours,humidity,temp,tempmax,tempmin,uvindex,solarradiation,solarenergy&unitGroup=metric", locationQuery, A.startDate, A.endDate, A.key)

	resp, err := makeHTTPRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		log.Printf("api fetch error: %s", err)
		log.Printf("error url: %s", url)
		return nil
	}

	defer resp.Body.Close()

	var apiResponse LocalMetrics
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		log.Printf("api response body decode error: %s\n", err)
		log.Printf("error url: %s", url)
		return nil
	}

	return &apiResponse
}

type probeAPIResponse struct {
	Locations map[string]struct {
		Address string `json:"address"`
	} `json:"locations"`
}

func (A *APIFetcher) ProbeCanonicalName(ctx context.Context, locationQuery string) (string, error) {
	log.Printf("probing canonical name for query %s\n", locationQuery)
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata/history?aggregateHours=24&startDateTime=2022-01-01T00:00:00&endDateTime=2022-01-02T00:00:00&unitGroup=metric&contentType=json&location=%s&key=%s", url2.QueryEscape(locationQuery), A.key)

	resp, err := makeHTTPRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", fmt.Errorf("canonical name probing error: %w", err)

	}

	defer resp.Body.Close()

	var apiResponse probeAPIResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	if err != nil {
		log.Printf("calling url: %s\n", url)
		return "", fmt.Errorf("canonical name probing error: %w", err)
	}

	if addressContainer, found := apiResponse.Locations[locationQuery]; found {
		if addressContainer.Address == "" {
			return "", fmt.Errorf("canonical name probing error: empty result")
		}
		return addressContainer.Address, nil
	} else {
		return "", fmt.Errorf("canonical name probing error: cannot find key '%s' in response %v", locationQuery, apiResponse.Locations)
	}
}

func makeHTTPRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}

func dumpResponseBody(response *http.Response) string {
	builder := new(strings.Builder)
	io.Copy(builder, response.Body)
	return builder.String()
}
