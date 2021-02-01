package ecdsa

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type GenerateEcdsaKeyFn func(ctx context.Context, req GenerateEcdsaKeyRequest) (response.Responser, error)

func NewGenerateEcdsaKeyFn() GenerateEcdsaKeyFn {
	return func(ctx context.Context, req GenerateEcdsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		privateKey, err := algorithm.GenerateEcdsaKey(ctx, req.CurveType)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		privatePem, err := algorithm.MarshalEcdsaPrivateKey(privateKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		publicPem, err := algorithm.MarshalEcdsaPublicKey(&privateKey.PublicKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		resp := GenerateEcdsaKeyResponse{
			PublicKey:  string(publicPem),
			PrivateKey: string(privatePem),
		}
		logger.Info("Create Public & Private Key Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &resp), nil
	}
}
