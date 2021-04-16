package rsa

import (
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type EncryptRsaKeyRequest struct {
	Message   string `json:"message" example:"trust"`
	PublicKey string `json:"publicKey" example:"-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAxH7IDHwozYyr4ZBvf0ySpc5XEDsYvXWmGEm7/bOQCp7m8NhjUXIV\n4AmPLQ3G0uz/1W10ZrzOWkUEC4LuVE7A4i2EY1qCR/F7UeVHq2/hycQrwq0QERgj\no+I2eUMBn5nXeP1s/rYWjhgUn0vp+VYNx/7e98UCO3hNWlinp01CVgtjCBZNH54H\n+nsmNYdQ63cppQJKpGHZ4TJgb9tb3dP2earUU0nCRSR/0+zdjIYlwTIaJAtNuxsT\nKcH/szXVuRMfM03CK/672FTK+5yzwup9EO349D9QZKl6GKOrqvtzsN0Sps+mppmD\nuEz4eAP3xFKZS+xN5+CwmE5ULGdZiHaQXQIDAQAB\n-----END RSA PUBLIC KEY-----\n"`
}

func (req *EncryptRsaKeyRequest) validate() error {
	if len(req.Message) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'message' must be REQUIRED field but the input is '%v'.", req.Message)), response.ValidateFieldError)
	}
	if len(req.PublicKey) == 0 {
		return errors.Wrap(errors.New(fmt.Sprintf("'publicKey' must be REQUIRED field but the input is '%v'.", req.PublicKey)), response.ValidateFieldError)
	}
	return nil
}

type EncryptRsaKeyResponse struct {
	EncryptedMessage string `json:"encryptedMessage" example:"MEUCIQCi"`
}
