package ecdsa

import (
	"context"
	"fmt"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type SignEcdsaKeyFn func(ctx context.Context, req SignEcdsaKeyRequest) (response.Responser, error)

func NewSignEcdsaKeyFn() SignEcdsaKeyFn {
	return func(ctx context.Context, req SignEcdsaKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		privateKey, err := algorithm.ParseEcdsaPrivateKeyFromPemStr(req.PrivateKey)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		hex, sign, err := algorithm.SignEcdsa(ctx, privateKey, req.HashType, req.Message)
		if err != nil {
			return response.NewErrResponse(response.ErrInternalServerCode, response.ErrBadRequestDesc, err.Error()), err
		}
		signEcdsaKeyResponse := SignEcdsaKeyResponse{
			SignHex:  fmt.Sprintf("%x", hex),
			SignByte: *sign,
		}
		logger.Info("Sign Message Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &signEcdsaKeyResponse), nil
	}
}
