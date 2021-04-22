package rsa

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type EncryptRsaKeyFn func(ctx context.Context, req EncryptRsaKeyRequest) (response.Responser, error)

func NewEncryptRsaKeyFn() EncryptRsaKeyFn {
	return func(ctx context.Context, req EncryptRsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		publicKey, err := algorithm.ParseRsaPublicKeyFromPemStr(req.PublicKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}
		encrypted, err := algorithm.EncryptRsa(ctx, publicKey, []byte(req.Message))
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}
		encryptRsaKeyResponse := EncryptRsaKeyResponse{
			EncryptedMessage: *encrypted,
		}
		logger.Info("Encrypted Message Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &encryptRsaKeyResponse), nil
	}
}
