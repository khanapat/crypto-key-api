package rsa

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type GenerateRsaKeyRequest struct {
	Bits int `json:"bits" example:"2048"`
}

func (req *GenerateRsaKeyRequest) validate() error {
	if req.Bits < 1 {
		return errors.Wrap(errors.New(fmt.Sprintf("'bits' must be REQUIRED field but the input is '%v'.", req.Bits)), response.ValidateFieldError)
	}
	return nil
}

type GenerateRsaKeyResponse struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
