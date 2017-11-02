package fohnhab

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	l "github.com/go-kit/kit/log"
)

// GenerateKeyRequest models an incoming request to the GenerateKeyEndpoint
type GenerateKeyRequest struct {
	Kind string `json:"kind"`
}

// GenerateKeyResponse models a response from the GenerateKeyEndpoint
type GenerateKeyResponse struct {
	Key string `json:"key"`
	Err string `json:"err,omitempty"`
}

// Endpoints models the collection of endpoints our service will use when being run
type Endpoints struct {
	GenerateKeyEndpoint endpoint.Endpoint
}

// Middleware serves as an endpoint decorater to implement logging at the transport and application level
type Middleware func(endpoint.Endpoint) endpoint.Endpoint

// MakeEndpoints returns the endpoints implemented by the service
func MakeEndpoints(svc Service, logger l.Logger) Endpoints {
	var ep Endpoints
	ep.GenerateKeyEndpoint = MakeGenerateKeyEndpoint(svc)
	ep.GenerateKeyEndpoint = transportMiddleware(l.With(logger, "method", "keygen"))(ep.GenerateKeyEndpoint)
	return ep
}

// MakeGenerateKeyEndpoint constructs an endpoint to be served byt our service
func MakeGenerateKeyEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req  GenerateKeyRequest
			resp GenerateKeyResponse
		)
		req = request.(GenerateKeyRequest)
		key, err := svc.GenerateKey(ctx, req)
		if err != nil {
			resp.Key = key
			resp.Err = err.Error()
			return resp, nil
		}
		resp.Key = key
		return resp, nil
	}
}

// DecodeGenerateKeyRequest converts and httpRequest to a request readable by our program
func DecodeGenerateKeyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GenerateKeyRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, err
}

// EncodeGenerateKeyResponse converts a given response struct to *http.Response to be written to the respective client
func EncodeGenerateKeyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := response.(GenerateKeyResponse).Err
	if err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
