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

// GCMERequest models an incoming request to the GCME..Endpoint
type GCMERequest struct {
	Key       string `json:"key"`
	ToEncrypt string `json:"plaintext"`
}

type GCMEResponse struct {
	CipherText string `json:"ciphertext"`
	Err        string `json:"err,omitempty"`
}

// GCMDRequest models an incoming request to the GCMDEncdpoint
type GCMDRequest struct {
	Key       string `json:"key"`
	ToDecrypt string `json:"ciphertext"`
}

type GCMDRespons struct {
	PlainText string `json:"plaintext"`
	Err       string `json:"err,omitempty"`
}

// Endpoints models the collection of endpoints our service will use when being run
type Endpoints struct {
	GenerateKeyEndpoint endpoint.Endpoint
	GCMEncryptEndpoint  endpoint.Endpoint
	GCMDecryptEndpoint  endpoint.Endpoint
}

// Middleware serves as an endpoint decorater to implement logging at the transport and application level
type Middleware func(endpoint.Endpoint) endpoint.Endpoint

// MakeEndpoints returns the endpoints implemented by the service
func MakeEndpoints(svc Service, logger l.Logger) Endpoints {
	var ep Endpoints
	ep.GenerateKeyEndpoint = MakeGenerateKeyEndpoint(svc)
	ep.GenerateKeyEndpoint = transportMiddleware(l.With(logger, "method", "keygen"))(ep.GenerateKeyEndpoint)
	ep.GCMEncryptEndpoint = MakeGMCEncryptEndpoint(svc)
	ep.GCMEncryptEndpoint = transportMiddleware(l.With(logger, "method", "encrypt"))(ep.GCMEncryptEndpoint)
	return ep
}

// MakeGenerateKeyEndpoint constructs an endpoint to be served by our service
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

// GenerateKey with an Endpoint implementation allows us to treat out EP struct as a Service
func (e Endpoints) GenerateKey(ctx context.Context, r GenerateKeyRequest) (interface{}, error) {
	resp, _ := e.GenerateKeyEndpoint(ctx, r)
	gkResp := resp.(GenerateKeyResponse)
	return gkResp, nil
}

// DecodeGenerateKeyRequest converts an httpRequest to a request readable by our program (Usable by server)
func DecodeGenerateKeyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GenerateKeyRequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, err
}

// EncodeGenerateKeyResponse converts a given response struct to *http.Response to be written to the respective client (Usable by server)
func EncodeGenerateKeyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := response.(GenerateKeyResponse).Err
	if err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// DecodeGenerateKeyResponse converts a given response to a json response (NEEDED FOR API GATEWAY IMPLEMENTATION)
func DecodeGenerateKeyResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var resp GenerateKeyResponse
	json.NewDecoder(r.Body).Decode(&resp)
	return resp, nil
}

func MakeGMCEncryptEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req  GCMERequest
			resp GCMEResponse
		)
		req = request.(GCMERequest)
		ct, err := svc.GCME(ctx, req)
		if err != nil {
			resp.CipherText = ct
			resp.Err = err.Error()
			return resp, nil
		}
		resp.CipherText = ct
		return resp, nil
	}
}

// GCME with an Endpoint receiver allows us to treat the ENDPOINT struct as a service
func (e Endpoints) GCME(ctx context.Context, r GCMERequest) (interface{}, error) {
	resp, _ := e.GCMEncryptEndpoint(ctx, r)
	GCMEresp := resp.(GCMEResponse)
	return GCMEresp, nil
}

func DecodeGCMEncryptRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GCMERequest
	var err error
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, err
}

func EncodeGCMEncryptResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := response.(GCMEResponse).Err
	if err != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Context-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
