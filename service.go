package fohnhab

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	l "github.com/go-kit/kit/log"
)

// Service is an interface that contains all of the endpoints the server will expose
type Service interface {
	GenerateKey(ctx context.Context, req GenerateKeyRequest) (string, error)
}

type fohnhabService struct{}

// NewService is a constructor for our fohnhab service
func NewService(logger l.Logger) Service {
	var svc Service
	svc = fohnhabService{}
	svc = logginMiddleware{logger, svc}
	count := configureRequestCount()
	hist := configureRequestLatency()
	svc = instrumentingMiddleware{count, hist, svc}
	return svc
}

func (fohnhabService) GenerateKey(ctx context.Context, req GenerateKeyRequest) (string, error) {
	var (
		key       []byte
		err       error
		kind      string
		keyString string
	)
	kind = req.Kind
	switch kind {
	case "aes-256":
		c := 32
		key = make([]byte, c)
		_, err := rand.Read(key)
		if err != nil {
			return "", err
		}
		keyString = base64.StdEncoding.EncodeToString(key)
	default:
		err = fmt.Errorf("Type %v not found", kind)
	}
	return keyString, err
}
