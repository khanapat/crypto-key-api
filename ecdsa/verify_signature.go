package ecdsa

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type VerifyEcdsaKeyFn func(ctx context.Context, req VerifyEcdsaKeyRequest) (response.Responser, error)

func NewVerifyEcdsaKeyFn() VerifyEcdsaKeyFn {
	return func(ctx context.Context, req VerifyEcdsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		publicKey, err := algorithm.ParseEcdsaPublicKeyFromPemStr(req.PublicKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		valid, err := algorithm.VerifyEcdsa(ctx, publicKey, req.HashType, req.Message, req.Signature)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		verifyEcdsaKeyResponse := VerifyEcdsaKeyResponse{
			Validation: valid,
		}
		logger.Info("Verify Signature Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &verifyEcdsaKeyResponse), nil
	}
}
