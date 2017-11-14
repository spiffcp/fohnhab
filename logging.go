package fohnhab

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Transport

// LoggingMiddleware accepts any logger implementation and returns a decorator that wraps an endpoint
func transportMiddleware(l log.Logger) Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			l.Log("msg", "calling endpoint")
			defer l.Log("msg", "called endpoint")
			return next(ctx, req)
		}
	}
}

// Application
type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) GenerateKey(ctx context.Context, req GenerateKeyRequest) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "generatekey",
			"input", req.Kind,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.GenerateKey(ctx, req)
	return
}

func (mw loggingMiddleware) GCME(ctx context.Context, req GCMERequest) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "gcm",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.GCME(ctx, req)
	return
}

func (mw loggingMiddleware) GCMD(ctx context.Context, req GCMDRequest) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "gcm",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.GCMD(ctx, req)
	return
}
