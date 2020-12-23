package key

import (
	"encoding/json"
	"net/http"

	"krungthai.com/khanapat/dpki/generate-key-api/middleware"
	"krungthai.com/khanapat/dpki/generate-key-api/response"
)

type asymmetricEcdsaKey struct {
	GenerateEcdsaKeyFn GenerateEcdsaKeyFn
}

func NewAsymmetricEcdsaKey(generateEcdsaKeyFn GenerateEcdsaKeyFn) http.Handler {
	return &asymmetricEcdsaKey{
		GenerateEcdsaKeyFn: generateEcdsaKeyFn,
	}
}

func (ak *asymmetricEcdsaKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := middleware.ContextData(r.Context())

	var req GenerateEcdsaKeyRequest
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

	resp, err := ak.GenerateEcdsaKeyFn(r.Context(), req)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&resp)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
