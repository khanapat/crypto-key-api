package key

import (
	"context"

	"krungthai.com/khanapat/dpki/crypto-key-api/algorithm"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type ValidatePublicKeyFn func(ctx context.Context, req ValidatePublicKeyRequest) (response.Responser, error)

func NewValidatePublicKeyFn() ValidatePublicKeyFn {
	return func(ctx context.Context, req ValidatePublicKeyRequest) (response.Responser, error) {
		logger := middleware.ContextData(ctx)

		var resp ValidatePublicKeyResponse
		publickeyType, err := algorithm.ValidatePublicKeyFromPemStr(req.PublicKey)
		if err != nil {
			return response.NewResponse(response.ErrInvalidRequestCode, response.ErrInvalidRequestDesc, &resp), err
		}
		resp.PublicKeyType = publickeyType
		resp.PublicKeyStatus = true
		logger.Info("Validate Key Success")
		return response.NewResponse(response.SuccessCode, response.SuccessDesc, &resp), nil
	}
}
