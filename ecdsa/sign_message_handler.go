package ecdsa

import (
	"encoding/json"
	"net/http"

	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type signEcdsaKey struct {
	SignEcdsaKeyFn SignEcdsaKeyFn
}

func NewSignEcdsaKey(signEcdsaKeyFn SignEcdsaKeyFn) http.Handler {
	return &signEcdsaKey{
		SignEcdsaKeyFn: signEcdsaKeyFn,
	}
}

// SignMessage example
// @Summary Sign Message
// @Description Method for signing signature.
// @Tags ECDSA
// @Accept json
// @Produce json
// @Param SignEcdsaKey body ecdsa.SignEcdsaKeyRequest true "object body for signing message."
// @Success 200 {object} response.Response{data=ecdsa.SignEcdsaKeyResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /ktb/blockchain/v1/crypto/ecdsa/sign [post]
func (s *signEcdsaKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ContextData(r.Context())

	var req SignEcdsaKeyRequest
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

	resp, err := s.SignEcdsaKeyFn(r.Context(), req)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
