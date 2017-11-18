package fohnhab

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	l "github.com/go-kit/kit/log"
)

// NonceSize - Set generic nonce size cited in RFC5084
const NonceSize = 12

// Service is an interface that contains all of the endpoints the server will expose
type Service interface {
	GenerateKey(ctx context.Context, req GenerateKeyRequest) (string, error)
	GCME(ctx context.Context, req GCMERequest) (string, error)
	GCMD(ctx context.Context, req GCMDRequest) (string, error)
}

type fohnhabService struct{}

// NewService is a constructor for our fohnhab service
func NewService(logger l.Logger) Service {
	var svc Service
	svc = fohnhabService{}
	svc = loggingMiddleware{logger, svc}
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
		rand.Read(key)
		keyString = base64.StdEncoding.EncodeToString(key)
	default:
		err = fmt.Errorf("Type %v not found", kind)
	}
	return keyString, err
}

func (fohnhabService) GCME(ctx context.Context, req GCMERequest) (string, error) {
	var (
		s string
		e error
	)
	key, _ := base64.StdEncoding.DecodeString(req.Key)
	plaintext := []byte(req.PlainText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, NonceSize)
	io.ReadFull(rand.Reader, nonce[:])
	out := gcm.Seal(nonce[:], nonce[:], plaintext, nil)
	s = base64.StdEncoding.EncodeToString(out)
	return s, e
}

func (fohnhabService) GCMD(ctx context.Context, req GCMDRequest) (string, error) {
	var (
		s string
		e error
	)
	key, _ := base64.StdEncoding.DecodeString(req.Key)
	ciphertext, _ := base64.StdEncoding.DecodeString(req.CipherText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, NonceSize)

	copy(nonce, ciphertext)
	plaintext, err := aesgcm.Open(nil, nonce[:], ciphertext[12:], nil)
	if err != nil {
		return "", err
	}
	s = string(plaintext[:])
	return s, e
}
