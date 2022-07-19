package rest

import (
	"encoding/json"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func testNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		type requestBody struct {
			Input []float64 `json:"input"`
		}
		var b requestBody
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			respondError(w, http.StatusBadRequest, ErrDecodingRequestBody)
			return
		}

		output, confidence, err := s.TestNetwork(ctx, networkId, b.Input)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "testing network"))
			return
		}

		type responseBody struct {
			Output     string  `json:"output"`
			Confidence float64 `json:"confidence"`
		}
		respondJson(w, http.StatusOK, responseBody{Output: output, Confidence: confidence})
		return
	}
}
