package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type TemperatureResponse struct {
	SensorID    string  `json:"sensorId"`
	Location    string  `json:"location"`
	Value float64 `json:"Value"`
	Unit        string  `json:"unit"`
	Ts          int64   `json:"ts"`
}

func resolveLocationAndSensor(location, sensorID string) (string, string) {
	location = strings.TrimSpace(location)
	sensorID = strings.TrimSpace(sensorID)

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
	return location, sensorID
}

func writeTemperature(w http.ResponseWriter, location, sensorID string) {
	location, sensorID = resolveLocationAndSensor(location, sensorID)
	rand.Seed(time.Now().UnixNano())
	t := -10.0 + rand.Float64()*(35.0-(-10.0))
	t = float64(int(t*10+0.5)) / 10.0

	resp := TemperatureResponse{
		SensorID:    sensorID,
		Location:    location,
		Value: t,
		Unit:        "C",
		Ts:          time.Now().Unix(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func temperatureAny(w http.ResponseWriter, r *http.Request, base string) {
	// Поддерживаем:
	// 1) GET base?location=&sensorId=
	// 2) GET base/{sensorId}
	path := strings.TrimPrefix(r.URL.Path, base)
	path = strings.Trim(path, "/")

	if path != "" {
		// вариант с path-параметром
		writeTemperature(w, "", path)
		return
	}
	// вариант с query-параметрами
	q := r.URL.Query()
	writeTemperature(w, q.Get("location"), q.Get("sensorId"))
}

func main() {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", healthz)

	// /temperature и все под-пути
	mux.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/temperature" && !strings.HasPrefix(r.URL.Path, "/temperature/") {
			http.NotFound(w, r)
			return
		}
		temperatureAny(w, r, "/temperature")
	})
	mux.HandleFunc("/temperature/", func(w http.ResponseWriter, r *http.Request) {
		temperatureAny(w, r, "/temperature")
	})

	// зеркала под /api/temperature и все под-пути (на случай, если смарт-дом зовёт именно так)
	mux.HandleFunc("/api/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/temperature" && !strings.HasPrefix(r.URL.Path, "/api/temperature/") {
			http.NotFound(w, r)
			return
		}
		temperatureAny(w, r, "/api/temperature")
	})
	mux.HandleFunc("/api/temperature/", func(w http.ResponseWriter, r *http.Request) {
		temperatureAny(w, r, "/api/temperature")
	})

	// корень
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("temperature-api-go: try /healthz or /temperature"))
			return
		}
		http.NotFound(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("temperature-api-go listening on :%s", port)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       75 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
