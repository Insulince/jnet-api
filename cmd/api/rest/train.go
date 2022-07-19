package rest

import (
	"encoding/json"
	"fmt"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/Insulince/jnet-api/untitled/cmd/api/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
)

func trainNetwork(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		type requestBody struct {
			LearningRate      float64 `json:"learningRate"`
			MiniBatchSize     int     `json:"miniBatchSize"`
			AverageLossCutoff float64 `json:"averageLossCutoff"`
			MinLossCutoff     float64 `json:"minLossCutoff"`
			MaxIterations     int     `json:"maxIterations"`
			// TODO(justin): Implement?
			// Timeout time.Duration `json:"timeout"`
		}
		var b requestBody
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			respondError(w, http.StatusBadRequest, ErrDecodingRequestBody)
			return
		}

		err := s.TrainNetwork(ctx, networkId, b.LearningRate, b.MiniBatchSize, b.AverageLossCutoff, b.MinLossCutoff, b.MaxIterations)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				respond(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			respondError(w, http.StatusBadRequest, errors.Wrap(err, "training network"))
			return
		}

		// TODO(justin): Is this right?
		respond(w, http.StatusAccepted, http.StatusText(http.StatusAccepted))
		return
	}
}

func trainingStatus(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		_, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func haltTraining(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		_, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func resetNetwork(_ service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Context()

		_, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		respond(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
		return
	}
}

func streamTrainingProgress(s service.Service, wsUpgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		networkId, found := mux.Vars(r)["networkId"]
		if !found {
			respondError(w, http.StatusBadRequest, ErrNoNetworkIdInPath)
			return
		}

		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			respondError(w, http.StatusInternalServerError, errors.Wrap(err, "upgrading connection"))
			return
		}

		// NOTE(justin): After upgrading the connection, you cannot send typical HTTP responses, so you much log instead.

		worker := ws.NewWorker(conn)
		defer worker.Close()
		fmt.Printf("new worker connected: %v\n", worker.Id)

		go func() {
			// NOTE(justin): Any message received is to be interpret as a signal to close the connection.
			_, _, _ = conn.ReadMessage()
			fmt.Printf("worker disconnected: %v\n", worker.Id)
			worker.ConnectionClosed <- struct{}{}
		}()

		err = s.StreamTrainingProgress(ctx, networkId, worker)
		if err != nil {
			fmt.Println(errors.Wrap(err, "streaming training progress"))
			return
		}
	}
}
