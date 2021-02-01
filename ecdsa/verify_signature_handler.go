package ecdsa

import (
	"encoding/json"
	"net/http"

	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type verifyEcdsaKey struct {
	VerifyEcdsaKeyFn VerifyEcdsaKeyFn
}

func NewVerifyEcdsaKey(verifyEcdsaKeyFn VerifyEcdsaKeyFn) http.Handler {
	return &verifyEcdsaKey{
		VerifyEcdsaKeyFn: verifyEcdsaKeyFn,
	}
}

// VerifySignature example
// @Summary Verify Signature
// @Description Method for verifying signature.
// @Tags ECDSA
// @Accept json
// @Produce json
// @Param VerifyEcdsaKey body ecdsa.VerifyEcdsaKeyRequest true "object body for verifying signature."
// @Success 200 {object} response.Response{data=ecdsa.VerifyEcdsaKeyResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /ktb/blockchain/v1/crypto/ecdsa/verify [post]
func (s *verifyEcdsaKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ContextData(r.Context())

	var req VerifyEcdsaKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrResponse(response.ErrBadRequestCode, response.ErrBadRequestDesc, nil))
		return
	}
	defer r.Body.Close()

	if err := req.validate(); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewErrResponse(response.ErrBadRequestCode, response.ErrBadRequestDesc, err.Error()))
		return
	}

	resp, err := s.VerifyEcdsaKeyFn(r.Context(), req)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
