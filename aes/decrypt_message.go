package aes

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type DecryptAesKeyFn func(ctx context.Context, req DecryptAesKeyRequest) (response.Responser, error)

func NewDecryptAesKeyFn() DecryptAesKeyFn {
	return func(ctx context.Context, req DecryptAesKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		plainText, err := algorithm.DecryptAesBlockGCM(ctx, req.Key, req.Salt, req.EncryptedText)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrInternalServerDesc, err.Error()), err
		}
		decryptAesKeyResponse := DecryptAesKeyResponse{
			DecryptedText: *plainText,
		}
		logger.Info("Decrypted Message Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &decryptAesKeyResponse), nil
	}
}
