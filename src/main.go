package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP recebidas.",
		},
		[]string{"endpoint", "method"},
	)
	serviceUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_server_up",
			Help: "Disponibilidade do serviço (1 = up, 0 = down).",
		},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(serviceUp)
}

type response struct {
	Name string `json:"nome"`
	Time string `json:"horario"`
}

func getData(w http.ResponseWriter, r *http.Request) {
	requestsTotal.WithLabelValues("/projeto-korp", r.Method).Inc()

	resp := response{
		Name: "Projeto Korp",
		Time: time.Now().String()[11:19],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))

}

func main() {

	serviceUp.Set(1)

	mux := http.NewServeMux()
	mux.HandleFunc("/projeto-korp", getData)
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
