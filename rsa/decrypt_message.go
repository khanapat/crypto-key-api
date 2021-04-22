package rsa

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type DecryptRsaKeyFn func(ctx context.Context, req DecryptRsaKeyRequest) (response.Responser, error)

func NewDecryptRsaKeyFn() DecryptRsaKeyFn {
	return func(ctx context.Context, req DecryptRsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		privateKey, err := algorithm.ParseRsaPrivateKeyFromPemStr(req.PrivateKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}

		decrypted, err := algorithm.DecryptRsa(ctx, privateKey, req.CipherText)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}
		decryptRsaKeyResponse := DecryptRsaKeyResponse{
			DecryptedMessage: *decrypted,
		}
		logger.Info("Decrypted Message Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &decryptRsaKeyResponse), nil
	}
}
