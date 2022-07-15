package main

import (
	"context"
)

type MetricFetcher interface {
	Fetch(ctx context.Context, locationQuery string) *LocalMetrics
}

type LocalMetrics struct {
	Address         string          `json:"address"`
	ResolvedAddress string          `json:"resolvedAddress"`
	Days            []*DailyMetrics `json:"days"`
}

type DailyMetrics struct {
	Date  string           `json:"datetime"`
	Hours []*HourlyMetrics `json:"hours"`
}

type HourlyMetrics struct {
	Hour           string   `json:"datetime"`
	Temp           float64  `json:"temp"`
	Humidity       float64  `json:"humidity"`
	SolarRadiation float64  `json:"solarradiation"`
	SolarEnergy    *float64 `json:"solarenergy"`
	UVIndex        float64  `json:"uvindex"`
}

type StubFetcher struct {
}

func (s *StubFetcher) Fetch(ctx context.Context, locationQuery string) *LocalMetrics {
	solarEnergy := 12.4
	return &LocalMetrics{
		ResolvedAddress: locationQuery,
		Days: []*DailyMetrics{
			&DailyMetrics{
				Date: "2021-01-01",
				Hours: []*HourlyMetrics{
					&HourlyMetrics{
						Hour:           "00:00:00",
						Temp:           24,
						Humidity:       42,
						SolarRadiation: 5,
						SolarEnergy:    &solarEnergy,
						UVIndex:        0.5,
					},
				},
			},
		},
	}
}
