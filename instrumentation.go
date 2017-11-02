package fohnhab

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var fieldKeys = []string{"method", "error"}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (mw instrumentingMiddleware) GenerateKey(ctx context.Context, req GenerateKeyRequest) (output string, err error) {
	defer func(begin time.Time) {
		instVal := []string{"method", "GenerateKey", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(instVal...).Add(1)
		mw.requestLatency.With(instVal...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.next.GenerateKey(ctx, req)
	return
}

func configureRequestCount() *kitprometheus.Counter {
	var c *kitprometheus.Counter
	c = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "meetup",
		Subsystem: "fohnhab_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	return c
}

func configureRequestLatency() *kitprometheus.Summary {
	var s *kitprometheus.Summary
	s = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "meetup",
		Subsystem: "fohnhab_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	return s
}
