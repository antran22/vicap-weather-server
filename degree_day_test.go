package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDegreeDay(t *testing.T) {
	sampleMetrics := loadSampleAPIResponse(t)
	degreeDay := CalculateDegreeDay(sampleMetrics.Days[0], 18)
	assert.InEpsilon(t, 6.6, degreeDay.Heating, 1e-6)
	assert.InEpsilon(t, 67.7, degreeDay.Cooling, 1e-6)
}
