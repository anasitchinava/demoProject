package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

var (
    RequestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "status"},
    )

    ResponseStatus = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "response_status",
            Help: "Status of HTTP responses",
        },
        []string{"status"},
    )
)

func Init() {
    prometheus.MustRegister(RequestCounter)
    prometheus.MustRegister(ResponseStatus)
}

func Handler() http.Handler {
    return promhttp.Handler()
}
