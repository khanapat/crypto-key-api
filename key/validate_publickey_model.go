package key

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type ValidatePublicKeyRequest struct {
	PublicKey string `json:"publicKey"`
}

func (req *ValidatePublicKeyRequest) validate() error {
	if len(req.PublicKey) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'publicKey' must be REQUIRED field but the input is '%v'.", req.PublicKey)), response.ValidateFieldError)
	}
	return nil
}

type ValidatePublicKeyResponse struct {
	PublicKeyType   string `json:"publicKeyType"`
	PublicKeyStatus bool   `json:"publicKeyStatus"`
}
