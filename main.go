package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func writeJSONBody(w http.ResponseWriter, body interface{}) {
	jsonOutput, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	if _, err = w.Write(jsonOutput); err != nil {
		log.Println(err)
	}
}

func main() {
	var fetcher MetricFetcher = NewCachedAPIFetcher()
	r := chi.NewRouter()

	r.Route("/{location}", func(r chi.Router) {
		r.Get("/degree_day", func(w http.ResponseWriter, r *http.Request) {
			baseTempString := r.URL.Query().Get("base_temp")
			if baseTempString == "" {
				w.WriteHeader(http.StatusBadRequest)
				writeJSONBody(w, map[string]string{
					"message": "query base_temp required",
				})
				return
			}
			baseTemp, err := strconv.ParseFloat(baseTempString, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				writeJSONBody(w, map[string]string{
					"message": "Cannot parse base_temp as a floating number",
				})
				return
			}

			location := chi.URLParam(r, "location")
			metrics := fetcher.Fetch(r.Context(), location)

			totalCoolingDegreeDay, totalHeatingDegreeDay := 0.0, 0.0

			for _, day := range metrics.Days {
				degreeDay := CalculateDegreeDay(day, baseTemp)
				totalHeatingDegreeDay += degreeDay.Heating
				totalCoolingDegreeDay += degreeDay.Cooling
			}

			writeJSONBody(w, map[string]float64{
				"heating_degree_day": totalHeatingDegreeDay,
				"cooling_degree_day": totalCoolingDegreeDay,
			})
		})
	})

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":8080"
	}

	log.Printf("Listening on %s\n", address)
	log.Fatalln(http.ListenAndServe(address, r))
}
