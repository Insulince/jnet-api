package rest

import (
	"encoding/json"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func addTrainingData(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		type requestBody struct {
			TrainingData []networks.TrainingDatum `json:"trainingData"`
		}
		var b requestBody
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			respondError(w, http.StatusBadRequest, ErrDecodingRequestBody)
			return
		}

		ids, err := s.AddTrainingData(ctx, networkId, b.TrainingData)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "adding training data"))
			return
		}

		type responseBody struct {
			Ids []string `json:"ids"`
		}
		respondJson(w, http.StatusCreated, responseBody{Ids: ids})
		return
	}
}

func searchTrainingData(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		tds, err := s.SearchTrainingData(ctx, networkId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "searching training data"))
		}

		type responseBody struct {
			TrainingData []networks.TrainingDatum `json:"trainingData"`
		}
		respondJson(w, http.StatusOK, responseBody{TrainingData: tds})
		return
	}
}

func deleteAllTrainingData(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		err := s.DeleteAllTrainingData(ctx, networkId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "deleting all training data"))
			return
		}

		respond(w, http.StatusNoContent, "")
		return
	}
}

func getTrainingDatum(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		trainingDatumId, found := mux.Vars(r)["trainingDatumId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoTrainingDatumIdInPath)
			return
		}

		td, err := s.GetTrainingDatum(ctx, trainingDatumId)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "getting training data"))
			return
		}

		respondJson(w, http.StatusOK, td)
		return
	}
}

func updateTrainingDatum(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		_, found := mux.Vars(r)["trainingDatumId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoTrainingDatumIdInPath)
			return
		}

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func patchTrainingDatum(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		_, found := mux.Vars(r)["trainingDatumId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoTrainingDatumIdInPath)
			return
		}

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func deleteTrainingDatum(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		trainingDatumId, found := mux.Vars(r)["trainingDatumId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoTrainingDatumIdInPath)
			return
		}

		err := s.DeleteTrainingDatum(ctx, trainingDatumId)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				respond(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "deleting training datum"))
			return
		}

		respond(w, http.StatusNoContent, "")
		return
	}
}
