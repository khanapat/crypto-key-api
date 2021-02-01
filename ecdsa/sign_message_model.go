package ecdsa

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type SignEcdsaKeyRequest struct {
	HashType   string `json:"hashType" example:"sha256"`
	Message    string `json:"message" example:"trust"`
	PrivateKey string `json:"privateKey" example:"-----BEGIN ECDSA PRIVATE KEY-----\nMHcCAQEEIB5g4Upn7ewh+vSLq9f4WJxdbhTfpsYa0SYaEkDl7xZPoAoGCCqGSM49\nAwEHoUQDQgAEmacWvMg72qXbAuh1JnfFwjY5eU1SxAiphgN3UQXTzlHJR0RGJsSL\nRuYpbc5asjL+oXvQ41ENxbYE58EsXMhbOw==\n-----END ECDSA PRIVATE KEY-----\n"`
}

func (req *SignEcdsaKeyRequest) validate() error {
	if len(req.HashType) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'hashType' must be REQUIRED field but the input is '%v'.", req.HashType)), response.ValidateFieldError)
	}
	if len(req.Message) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'message' must be REQUIRED field but the input is '%v'.", req.Message)), response.ValidateFieldError)
	}
	if len(req.PrivateKey) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'privateKey' must be REQUIRED field but the input is '%v'.", req.PrivateKey)), response.ValidateFieldError)
	}
	return nil
}

type SignEcdsaKeyResponse struct {
	SignHex  string `json:"signHex" example:"3045022100a2f880d14e5fdc53596cbb904b0ea9c2418d75a32d9fb107401966c585ea1d1d022001cf0e743e5a853da97864f5476733a3d4ec4e00972921b28d40efa7769d15a2"`
	SignByte string `json:"signByte" example:"MEUCIQCi+IDRTl/cU1lsu5BLDqnCQY11oy2fsQdAGWbFheodHQIgAc8OdD5ahT2peGT1R2czo9TsTgCXKSGyjUDvp3adFaI="`
}
