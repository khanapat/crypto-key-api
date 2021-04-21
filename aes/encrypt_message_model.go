package aes

import (
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type EncryptAesKeyRequest struct {
	Key     string `json:"key" example:"passphrasewhichneedstobe32bytes!"`
	Payload string `json:"payload" example:"myfile"`
}

func (req *EncryptAesKeyRequest) validate() error {
	if utf8.RuneCountInString(req.Key) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'key' must be REQUIRED field but the input is '%v'.", req.Key)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Payload) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'payload' must be REQUIRED field but the input is '%v'.", req.Payload)), response.ValidateFieldError)
	}
	return nil
}

type EncryptAesKeyResponse struct {
	CipherText string `json:"cipherText" example:"d2edb6726d74e4c0bdf9a3211e0be0c93969fb7da4aad3a916ab7e5e1d8e09c6"`
	Salt       string `json:"salt" example:"2abfd6e910b10f4bf07a861e"`
}
