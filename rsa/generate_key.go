package rsa

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type GenerateRsaKeyFn func(ctx context.Context, req GenerateRsaKeyRequest) (response.Responser, error)

func NewGenerateRsaKeyFn() GenerateRsaKeyFn {
	return func(ctx context.Context, req GenerateRsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		privateKey, err := algorithm.GenerateRsaKey(ctx, req.Bits)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		privatePem, err := algorithm.MarshalRsaPrivateKey(privateKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		publicPem, err := algorithm.MarshalRsaPublicKey(&privateKey.PublicKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		resp := GenerateRsaKeyResponse{
			PublicKey:  string(publicPem),
			PrivateKey: string(privatePem),
		}
		logger.Info("Create Public & Private Key Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &resp), nil
	}
}
