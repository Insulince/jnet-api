package service

import (
	"context"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet-api/untitled/cmd/api/ws"
	"github.com/Insulince/jnet/pkg/trainer"
	"github.com/pkg/errors"
)

func (s Service) TrainNetwork(ctx context.Context, networkId string, lr float64, mbs int, alc float64, mlc float64, mi int) error {
	// TODO(justin): make these steps concurrent
	nw, _, err := s.GetNetworkById(ctx, networkId, false)
	if err != nil {
		return errors.Wrap(err, "getting network by id")
	}
	tds, err := s.SearchTrainingData(ctx, networkId)
	if err != nil {
		return errors.Wrap(err, "searching training data")
	}

	jnw, err := s.networkTranslator.Deserialize(nw.Blueprint)
	if err != nil {
		return errors.Wrap(err, "deserializing network")
	}

	var jtds trainer.Data
	for _, td := range tds {
		var jtd trainer.Datum
		jtd.Data = td.Data
		jtd.Truth = td.Truth
		jtds = append(jtds, jtd)
	}

	var tc trainer.Configuration
	tc.LearningRate = lr
	tc.MiniBatchSize = mbs
	tc.AverageLossCutoff = alc
	tc.MinLossCutoff = mlc
	tc.MaxIterations = mi
	// TODO(justin): Implement?
	// tc.Timeout = 0

	var ts networks.TrainingSpec
	ts.NetworkId = networkId
	ts.Network = jnw
	ts.Data = jtds
	ts.TrainingConfiguration = tc

	go func() { s.trainingQueue <- ts }()

	return nil
}

func (s Service) StreamTrainingProgress(_ context.Context, networkId string, worker *ws.Worker) error {
	for {
		select {
		case message := <-worker.Stream:
			if message.NetworkId != networkId {
				continue
			}
			if err := worker.SendMessage(message.Data); err != nil {
				return errors.Wrap(err, "sending message")
			}
		case <-worker.ConnectionClosed:
			return nil
		}
	}
}
