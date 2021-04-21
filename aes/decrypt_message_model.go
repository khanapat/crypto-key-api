package aes

import (
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type DecryptAesKeyRequest struct {
	Key           string `json:"key" example:"passphrasewhichneedstobe32bytes!"`
	Salt          string `json:"salt" example:"2abfd6e910b10f4bf07a861e"`
	EncryptedText string `json:"encryptedText" example:"Zrecndea6xeD6WDPPy9Tau0drlty"`
}

func (req *DecryptAesKeyRequest) validate() error {
	if utf8.RuneCountInString(req.Key) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'key' must be REQUIRED field but the input is '%v'.", req.Key)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.Salt) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'salt' must be REQUIRED field but the input is '%v'.", req.Salt)), response.ValidateFieldError)
	}
	if utf8.RuneCountInString(req.EncryptedText) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'encryptedText' must be REQUIRED field but the input is '%v'.", req.EncryptedText)), response.ValidateFieldError)
	}
	return nil
}

type DecryptAesKeyResponse struct {
	DecryptedText string `json:"decryptedText" example:"trust"`
}
