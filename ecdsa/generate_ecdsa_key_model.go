package ecdsa

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/constant"
)

type GenerateEcdsaKeyRequest struct {
	CurveType string `json:"curveType"`
}

func (req *GenerateEcdsaKeyRequest) validate() error {
	if len(req.CurveType) <= 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'curveType' must be REQUIRED field but the input is '%v'.", req.CurveType)), constant.ValidateFieldError)
	}
	return nil
}

type GenerateEcdsaKeyResponse struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
