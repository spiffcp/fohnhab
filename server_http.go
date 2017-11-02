package fohnhab

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// NewHTTPServer returns a server construct that contains current endpoints along with their decode and encode methods
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/keygen", httptransport.NewServer(
		endpoints.GenerateKeyEndpoint,
		DecodeGenerateKeyRequest,
		EncodeResponse,
	))
	return m
}
