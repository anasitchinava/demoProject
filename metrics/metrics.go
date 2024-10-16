package metrics

import (
    "log"
    // "net/http"
    "sync"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/gin-gonic/gin"
    "strconv"
)

var (
    RequestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Number of get requests.",
        },
        []string{"path"},
    )

    ResponseStatus = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "response_status",
            Help: "Status of HTTP responses",
        },
        []string{"status"},
    )

    ResponseDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "http_response_time_seconds",
        Help:    "Duration of HTTP requests.",
        // Buckets: prometheus.DefBuckets,
    }, []string{"path"})

    mu          sync.Mutex
    initialized bool
)

func Init() {
    mu.Lock()
    defer mu.Unlock()

    if initialized {
        return
    }

    if err := prometheus.Register(RequestCounter); err != nil {
        log.Printf("Error registering RequestCounter: %v", err)
    }
    if err := prometheus.Register(ResponseStatus); err != nil {
        log.Printf("Error registering ResponseStatus: %v", err)
    }
    if err := prometheus.Register(ResponseDuration); err != nil {
        log.Printf("Error registering ResponseDuration: %v", err)
    }

    initialized = true
}

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Seconds()
        path := c.Request.URL.Path
        // method := c.Request.Method

        ResponseDuration.WithLabelValues(path).Observe(duration)

        status := c.Writer.Status()

        RequestCounter.WithLabelValues(path).Inc()
        ResponseStatus.WithLabelValues(strconv.Itoa(status)).Inc()
    }
}

// return the Prometheus handler for metrics
func Handler() gin.HandlerFunc {
    return gin.WrapH(promhttp.Handler())
}
