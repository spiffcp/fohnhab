package fohnhab

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	l "github.com/go-kit/kit/log"
)

// https://tools.ietf.org/html/rfc5084
// A nonce value of 12 octets can be processed more efficiently, so that length is RECOMMENDED.
// http://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf
// Designed by McGrew and Viega (GCM)
// AES designed by Joan Daemen, Vincent Rijmen
// Maximum Encrypted Plaintext Size: ≤ 239 – 256 bits

// NonceSize - Set generic nonce size cited in RFC5084
const NonceSize = 12

// Service is an interface that contains all of the endpoints the fohnhab server will expose
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

// Generate key provides key generation for different types of symetric and asymetric keys
func (fohnhabService) GenerateKey(ctx context.Context, req GenerateKeyRequest) (string, error) {
	var (
		key       []byte
		err       error
		kind      string
		keyString string
		c         int
	)
	kind = req.Kind
	switch kind {
	case "aes-256":
		c = 32
	case "aes-192":
		c = 24
	case "aes-128":
		c = 16
	default:
		err = fmt.Errorf("Type %v not found", kind)
	}
	key = make([]byte, c)
	rand.Read(key)
	keyString = base64.StdEncoding.EncodeToString(key)
	return keyString, err
}

// GCME accepts a defined request with given plaintext and key. It returns ciphertext.
func (fohnhabService) GCME(ctx context.Context, req GCMERequest) (string, error) {
	var (
		s string
		e error
	)

	key, err := base64.StdEncoding.DecodeString(req.Key)
	if err != nil {
		return "", errBase64DecodeKey
	}

	plaintext := []byte(req.PlainText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// NewGCM will error if the block is invalid - since we have already sanitized this input
	// We don't need to do any error handeling
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, NonceSize)
	io.ReadFull(rand.Reader, nonce[:])

	out := gcm.Seal(nil, nonce[:], plaintext, nil)
	s = base64.StdEncoding.EncodeToString(out)
	return s, e
}

// GCMD accepts a defined request with given ciphertext and key. It returns the plaintext.
func (fohnhabService) GCMD(ctx context.Context, req GCMDRequest) (string, error) {
	var (
		s string
		e error
	)

	key, err := base64.StdEncoding.DecodeString(req.Key)
	if err != nil {
		return "", errBase64DecodeKey
	}

	ciphertext, err := base64.StdEncoding.DecodeString(req.CipherText)
	if err != nil {
		return "", errBase64ErrDecodeCipher
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errBlockError
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

var errBlockError = errors.New("text is shorter than minumum AES block")
var errBase64DecodeKey = errors.New("unable to decode supplied key - check to make sure it is base64 encoded")
var errBase64ErrDecodeCipher = errors.New("unable to decode supplied ciphertext - check to make sure it is base64 encoded")
