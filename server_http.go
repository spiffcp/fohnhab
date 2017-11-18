package fohnhab

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	httptransport "github.com/go-kit/kit/transport/http"
)

// NewHTTPServer returns a server construct that contains current endpoints along with their decode and encode methods
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/keygen", httptransport.NewServer(
		endpoints.GenerateKeyEndpoint,
		DecodeGenerateKeyRequest,
		EncodeGenerateKeyResponse,
	))
	m.Handle("/encrypt", httptransport.NewServer(
		endpoints.GCMEncryptEndpoint,
		DecodeGCMEncryptRequest,
		EncodeGCMEncryptResponse,
	))
	m.Handle("/decrypt", httptransport.NewServer(
		endpoints.GCMDecryptEndpoint,
		DecodeGCMDecryptRequest,
		EncodeGCMDecryptResponse,
	))
	m.Handle("/metrics", promhttp.Handler())
	return m
}
