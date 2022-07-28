package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrDecodingRequestBody     = errors.New("decoding request body")
	ErrNoNetworkIdInPath       = errors.New("no network id in path")
	ErrNoTrainingDatumIdInPath = errors.New("no training datum id in path")
)

func NewRouter(s service.Service, wsUpgrader websocket.Upgrader) *mux.Router {
	r := mux.NewRouter()

	// TODO(justin): Create network FROM blueprint
	// NETWORKS
	r.HandleFunc("/networks", createNetwork(s)).Methods(http.MethodPost)
	r.HandleFunc("/networks", searchNetworks(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}", getNetwork(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}", updateNetwork(s)).Methods(http.MethodPut)
	r.HandleFunc("/networks/{networkId}", patchNetwork(s)).Methods(http.MethodPatch)
	r.HandleFunc("/networks/{networkId}", deleteNetwork(s)).Methods(http.MethodDelete)
	r.HandleFunc("/networks/{networkId}/activate", activateNetwork(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}/deactivate", deactivateNetwork(s)).Methods(http.MethodGet)

	// TRAINING DATA
	r.HandleFunc("/networks/{networkId}/training-data", addTrainingData(s)).Methods(http.MethodPost)
	r.HandleFunc("/networks/{networkId}/training-data", searchTrainingData(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}/training-data", deleteAllTrainingData(s)).Methods(http.MethodDelete)
	r.HandleFunc("/networks/-/training-data/{trainingDatumId}", getTrainingDatum(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/-/training-data/{trainingDatumId}", updateTrainingDatum(s)).Methods(http.MethodPut)
	r.HandleFunc("/networks/-/training-data/{trainingDatumId}", patchTrainingDatum(s)).Methods(http.MethodPatch)
	r.HandleFunc("/networks/-/training-data/{trainingDatumId}", deleteTrainingDatum(s)).Methods(http.MethodDelete)

	// TRAIN
	r.HandleFunc("/networks/{networkId}/train", trainNetwork(s)).Methods(http.MethodPost)
	r.HandleFunc("/networks/{networkId}/train/status", trainingStatus(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}/train/halt", haltTraining(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}/train/reset", resetNetwork(s)).Methods(http.MethodGet)
	r.HandleFunc("/networks/{networkId}/train/socket", streamTrainingProgress(s, wsUpgrader)).Methods(http.MethodGet)

	// TEST
	r.HandleFunc("/networks/{networkId}/test", testNetwork(s)).Methods(http.MethodPost)

	// HEALTH CHECK
	r.HandleFunc("/health", healthCheck(s)).Methods(http.MethodGet)

	r.Use(headerMiddleware)

	return r
}

func healthCheck(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		respond(w, http.StatusOK, http.StatusText(http.StatusOK))
		return
	}
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		next.ServeHTTP(w, r)
	})
}

func respondError(w http.ResponseWriter, status int, err error) {
	fmt.Println(err.Error())
	w.WriteHeader(status)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(err.Error())))
	if _, err := w.Write([]byte(err.Error())); err != nil {
		fmt.Println(errors.Wrap(err, "writing response body"))
		return
	}
}

func respond(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	if _, err := w.Write([]byte(body)); err != nil {
		fmt.Println(errors.Wrap(err, "writing response body"))
		return
	}
}

func respondJson(w http.ResponseWriter, status int, body interface{}) {
	js, err := json.Marshal(body)
	if err != nil {
		fmt.Println(errors.Wrap(err, "marshaling body into json"))
		respond(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	respond(w, status, string(js))
	return
}
