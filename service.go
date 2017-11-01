package fohnhab

import (
	"crypto/rand"
	"fmt"
)

// Service is an interface that contains all of the endpoints the server will expose
type Service interface {
	GenerateKey(kind string) ([]byte, error)
}

type fohnhabService struct{}

// NewService is a constructor for our fohnhab service
func NewService() Service {
	return fohnhabService{}
}

func (fohnhabService) GenerateKey(kind string) ([]byte, error) {
	var (
		key []byte
		err error
	)
	switch kind {
	case "aes-256":
		c := 32
		key = make([]byte, c)
		_, err := rand.Read(key)
		if err != nil {
			return nil, err
		}
	default:
		err = fmt.Errorf("Type %v not found", kind)
	}
	return key, err
}
