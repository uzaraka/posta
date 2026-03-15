/*
 *  MIT License
 *
 * Copyright (c) 2026 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jkaninda/okapi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "posta_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "posta_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	emailsSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "posta_emails_sent_total",
			Help: "Total number of emails sent successfully",
		},
	)

	emailsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "posta_emails_failed_total",
			Help: "Total number of emails that failed to send",
		},
	)

	emailsQueuedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "posta_emails_queued_total",
			Help: "Total number of emails enqueued for delivery",
		},
	)

	emailRetriesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "posta_email_retries_total",
			Help: "Total number of email retry attempts",
		},
	)

	webhookDeliveriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "posta_webhook_deliveries_total",
			Help: "Total number of webhook delivery attempts by status",
		},
		[]string{"status"},
	)

	webhookDeliveryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "posta_webhook_delivery_duration_seconds",
			Help:    "Duration of webhook delivery attempts in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	bouncesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "posta_bounces_total",
			Help: "Total number of bounces recorded by type",
		},
		[]string{"type"},
	)

	suppressionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "posta_suppressions_total",
			Help: "Total number of email suppressions added",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(emailsSentTotal)
	prometheus.MustRegister(emailsFailedTotal)
	prometheus.MustRegister(emailsQueuedTotal)
	prometheus.MustRegister(emailRetriesTotal)
	prometheus.MustRegister(webhookDeliveriesTotal)
	prometheus.MustRegister(webhookDeliveryDuration)
	prometheus.MustRegister(bouncesTotal)
	prometheus.MustRegister(suppressionsTotal)
}

// IncrementEmailSent increments the emails sent counter.
func IncrementEmailSent() {
	emailsSentTotal.Inc()
}

// IncrementEmailFailed increments the emails failed counter.
func IncrementEmailFailed() {
	emailsFailedTotal.Inc()
}

// IncrementEmailQueued increments the emails queued counter.
func IncrementEmailQueued() {
	emailsQueuedTotal.Inc()
}

// IncrementEmailRetry increments the email retries counter.
func IncrementEmailRetry() {
	emailRetriesTotal.Inc()
}

// IncrementWebhookDelivery increments the webhook delivery counter for the given status.
func IncrementWebhookDelivery(status string) {
	webhookDeliveriesTotal.WithLabelValues(status).Inc()
}

// ObserveWebhookDeliveryDuration records a webhook delivery duration.
func ObserveWebhookDeliveryDuration(seconds float64) {
	webhookDeliveryDuration.Observe(seconds)
}

// IncrementBounce increments the bounce counter for the given type (hard, soft, complaint).
func IncrementBounce(bounceType string) {
	bouncesTotal.WithLabelValues(bounceType).Inc()
}

// IncrementSuppression increments the suppression counter.
func IncrementSuppression() {
	suppressionsTotal.Inc()
}

// PrometheusMiddleware records HTTP request metrics.
func PrometheusMiddleware() okapi.Middleware {
	return func(c *okapi.Context) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		method := c.Request().Method
		path := c.Path()
		status := strconv.Itoa(c.Response().StatusCode())

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}

}

// MetricsHandler returns the Prometheus metrics handler.
func MetricsHandler() okapi.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *okapi.Context) error {
		handler.ServeHTTP(c.Response().(http.ResponseWriter), c.Request())
		return nil
	}
}
