package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TemperatureResponse struct {
	Value float64 `json:"value"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	http.HandleFunc("/temperature/", getTemperatureBySensor)
	http.HandleFunc("/temperature", getTemperatureByLocation)

	port := 8081
	log.Printf("🌡️  Temperature API running on port %d...", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getTemperatureByLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, `{"error":"missing parameter: location"}`, http.StatusBadRequest)
		return
	}

	temp := generateTemperature(location)
	log.Printf("→ location=%s → temperature=%.2f°C", location, temp)
	sendJSON(w, TemperatureResponse{Value: temp})
}

func getTemperatureBySensor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, `{"error":"missing sensorId"}`, http.StatusBadRequest)
		return
	}
	sensorID := parts[1]

	location := mapSensorToLocation(sensorID)
	temp := generateTemperature(location)
	log.Printf("→ sensorId=%s (%s) → temperature=%.2f°C", sensorID, location, temp)

	sendJSON(w, TemperatureResponse{Value: temp})
}

func generateTemperature(location string) float64 {
	base := map[string]float64{
		"Living Room": 22.0,
		"Bedroom":     20.0,
		"Kitchen":     25.0,
		"Unknown":     21.0,
	}
	temp := base[location] + rand.Float64()*4 - 2
	return math.Round(temp*100) / 100
}

func mapSensorToLocation(sensorID string) string {
	switch sensorID {
	case "1":
		return "Living Room"
	case "2":
		return "Bedroom"
	case "3":
		return "Kitchen"
	default:
		return "Unknown"
	}
}

func sendJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
