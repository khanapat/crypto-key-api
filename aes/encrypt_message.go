package aes

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type EncryptAesKeyFn func(ctx context.Context, req EncryptAesKeyRequest) (response.Responser, error)

func NewEncryptAesKeyFn() EncryptAesKeyFn {
	return func(ctx context.Context, req EncryptAesKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		cipherText, salt, err := algorithm.EncryptAesBlockGCM(ctx, req.Key, req.Payload)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}
		encryptAesKeyResponse := EncryptAesKeyResponse{
			CipherText: *cipherText,
			Salt:       *salt,
		}
		logger.Info("Encrypted Message Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &encryptAesKeyResponse), nil
	}
}
