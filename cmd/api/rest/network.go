package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

const (
	HeaderHideNetworkBlueprints = "X-JNet-Hide-Network-Blueprints"
	HeaderNetworkFromCache      = "X-JNet-Network-From-Cache"
)

func createNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		type requestBody struct {
			NeuronMap          []int    `json:"neuronMap"`
			InputLabels        []string `json:"inputLabels"`
			OutputLabels       []string `json:"outputLabels"`
			ActivationFunction string   `json:"activationFunction"`
		}
		var b requestBody
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			respondError(w, http.StatusBadRequest, ErrDecodingRequestBody)
			return
		}

		id, err := s.CreateNetwork(ctx, b.NeuronMap, b.InputLabels, b.OutputLabels, b.ActivationFunction)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "creating network"))
			return
		}

		w.Header().Set("Location", id)
		respond(w, http.StatusCreated, http.StatusText(http.StatusCreated))
		return
	}
}

func searchNetworks(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		nws, err := s.SearchNetworks(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "searching networks"))
			return
		}

		hideBlueprints := r.Header.Get(HeaderHideNetworkBlueprints)
		if hideBlueprints != "" {
			for i := range nws {
				nws[i].Blueprint = nil
			}
		}

		type responseBody struct {
			Networks []networks.Network `json:"networks"`
		}
		respondJson(w, http.StatusOK, responseBody{Networks: nws})
		return
	}
}

func getNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		nw, fromCache, err := s.GetNetworkById(ctx, networkId, true)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				respond(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "getting network by id"))
			return
		}

		w.Header().Set(HeaderNetworkFromCache, fmt.Sprint("%v", fromCache))
		respondJson(w, http.StatusOK, nw)
		return
	}
}

func updateNetwork(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func patchNetwork(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func deleteNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		err := s.DeleteNetwork(ctx, networkId)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				respond(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "deleting network"))
			return
		}

		respond(w, http.StatusNoContent, "")
		return
	}
}

func activateNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		err := s.ActivateNetwork(ctx, networkId)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				respond(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "activating network"))
			return
		}

		respond(w, http.StatusNoContent, "")
		return
	}
}

func deactivateNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		err := s.DeactivateNetwork(ctx, networkId)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "deactivating network"))
			return
		}

		respond(w, http.StatusNoContent, "")
		return
	}
}
