package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	promMetrics "github.com/kwa0x2/AutoSRT-Backend/monitoring/prometheus"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		promMetrics.ActiveConnections.Inc()

		c.Next()

		duration := time.Since(start).Seconds()
		method := c.Request.Method
		endpoint := c.FullPath()
		statusCode := strconv.Itoa(c.Writer.Status())

		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		promMetrics.HttpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()

		promMetrics.HttpRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)
		if c.Writer.Status() >= 500 {
			promMetrics.HttpErrorsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
		}

		promMetrics.ActiveConnections.Dec()
	}
}

func RecordSRTMetrics(status string, duration time.Duration) {
	if status == "queued_success" {
		promMetrics.QuededSRTRequest.WithLabelValues(status).Inc()
	}
}
