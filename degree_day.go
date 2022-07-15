package main

type DegreeDay struct {
	Heating float64
	Cooling float64
}

func CalculateDegreeDay(dailyMetrics *DailyMetrics, baseTemperature float64) DegreeDay {
	coldHourDiffAcc := 0.0
	hotHourDiffAcc := 0.0
	for _, hour := range dailyMetrics.Hours {
		if hour.Temp > baseTemperature {
			hotHourDiffAcc += hour.Temp - baseTemperature
		} else {
			coldHourDiffAcc += baseTemperature - hour.Temp
		}
	}

	return DegreeDay{
		Heating: coldHourDiffAcc,
		Cooling: hotHourDiffAcc,
	}
}
