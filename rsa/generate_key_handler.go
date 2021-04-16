package rsa

import (
	"encoding/json"
	"net/http"

	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type asymmetricRsaKey struct {
	GenerateRsaKeyFn GenerateRsaKeyFn
}

func NewAsymmetricRsaKey(generateRsaKeyFn GenerateRsaKeyFn) http.Handler {
	return &asymmetricRsaKey{
		GenerateRsaKeyFn: generateRsaKeyFn,
	}
}

// GenerateKey example
// @Summary Generate RSA Key
// @Description Method for generating key.
// @Tags RSA
// @Accept json
// @Produce json
// @Param GenerateRsaKey body rsa.GenerateRsaKeyRequest true "object body for generating key."
// @Success 200 {object} response.Response{data=rsa.GenerateRsaKeyResponse} "Success"
// @Failure 400 {object} response.ErrResponse "Bad Request"
// @Failure 500 {object} response.ErrResponse "Internal Server Error"
// @Router /ktb/blockchain/v1/crypto/rsa [post]
func (ak *asymmetricRsaKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ContextData(r.Context())

	var req GenerateRsaKeyRequest
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

	resp, err := ak.GenerateRsaKeyFn(r.Context(), req)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
