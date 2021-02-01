package ecdsa

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type VerifyEcdsaKeyRequest struct {
	HashType  string `json:"hashType" example:"sha256"`
	Message   string `json:"message" example:"trust"`
	Signature string `json:"signature" example:"MEUCIQCi+IDRTl/cU1lsu5BLDqnCQY11oy2fsQdAGWbFheodHQIgAc8OdD5ahT2peGT1R2czo9TsTgCXKSGyjUDvp3adFaI="`
	PublicKey string `json:"publicKey" example:"-----BEGIN ECDSA PRIVATE KEY-----\nMHcCAQEEIB5g4Upn7ewh+vSLq9f4WJxdbhTfpsYa0SYaEkDl7xZPoAoGCCqGSM49\nAwEHoUQDQgAEmacWvMg72qXbAuh1JnfFwjY5eU1SxAiphgN3UQXTzlHJR0RGJsSL\nRuYpbc5asjL+oXvQ41ENxbYE58EsXMhbOw==\n-----END ECDSA PRIVATE KEY-----\n"`
}

func (req *VerifyEcdsaKeyRequest) validate() error {
	if len(req.HashType) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'hashType' must be REQUIRED field but the input is '%v'.", req.HashType)), response.ValidateFieldError)
	}
	if len(req.Message) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'message' must be REQUIRED field but the input is '%v'.", req.Message)), response.ValidateFieldError)
	}
	if len(req.Signature) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'signature' must be REQUIRED field but the input is '%v'.", req.Signature)), response.ValidateFieldError)
	}
	if len(req.PublicKey) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'publicKey' must be REQUIRED field but the input is '%v'.", req.PublicKey)), response.ValidateFieldError)
	}
	return nil
}

type VerifyEcdsaKeyResponse struct {
	Validation bool `json:"validation" example:"true"`
}
