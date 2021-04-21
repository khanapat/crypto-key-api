package aes

import (
	"encoding/json"
	"net/http"

	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
	"krungthai.com/khanapat/dpki/crypto-key-api/response"
)

type decryptAesKey struct {
	DecryptAesKeyFn DecryptAesKeyFn
}

func NewDecryptAesKey(decryptAesKeyFn DecryptAesKeyFn) http.Handler {
	return &decryptAesKey{
		DecryptAesKeyFn: decryptAesKeyFn,
	}
}

func (s *decryptAesKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ContextData(r.Context())

	var req DecryptAesKeyRequest
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

	resp, err := s.DecryptAesKeyFn(r.Context(), req)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
