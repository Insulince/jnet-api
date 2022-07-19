package service

import (
	"context"
	"github.com/pkg/errors"
)

func (s Service) TestNetwork(ctx context.Context, networkId string, input []float64) (string, float64, error) {
	nw, fromCache, err := s.GetNetworkById(ctx, networkId, true)
	if err != nil {
		return "", 0, errors.Wrap(err, "getting network by id")
	}

	if !fromCache {
		err = s.ActivateNetwork(ctx, networkId)
		if err != nil {
			return "", 0, errors.Wrap(err, "activating network")
		}
	}

	output, confidence, err := nw.Network.Predict(input)
	if err != nil {
		return "", 0, errors.Wrap(err, "network prediction")
	}

	return output, confidence, nil
}
