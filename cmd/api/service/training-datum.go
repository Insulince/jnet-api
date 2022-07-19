package service

import (
	"context"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s Service) AddTrainingData(ctx context.Context, networkId string, trainingData []networks.TrainingDatum) ([]string, error) {
	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return nil, errors.Wrap(err, "object id from hex")
	}

	var tds []networks.TrainingDatum
	for _, jtd := range trainingData {
		var td networks.TrainingDatum
		td.NetworkId = pNetworkId
		td.Data = jtd.Data
		td.Truth = jtd.Truth
		tds = append(tds, td)
	}

	pids, err := s.mongo.InsertTrainingData(ctx, tds)
	if err != nil {
		if errors.Is(err, mongo.ErrEmptySlice) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "inserting training data")
	}

	var ids []string
	for _, pid := range pids {
		id := pid.Hex()
		ids = append(ids, id)
	}

	return ids, nil
}

func (s Service) SearchTrainingData(ctx context.Context, networkId string) ([]networks.TrainingDatum, error) {
	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return nil, errors.Wrap(err, "object id from hex")
	}

	tds, err := s.mongo.SearchTrainingData(ctx, pNetworkId)
	if err != nil {
		return nil, errors.Wrap(err, "searching training data")
	}

	return tds, nil
}

func (s Service) DeleteAllTrainingData(ctx context.Context, networkId string) error {
	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return errors.Wrap(err, "object id from hex")
	}

	err = s.mongo.DeleteAllTrainingData(ctx, pNetworkId)
	if err != nil {
		return errors.Wrap(err, "deleting all training data")
	}

	return nil
}

func (s Service) GetTrainingDatum(ctx context.Context, trainingDatumId string) (networks.TrainingDatum, error) {
	pTrainingDatumId, err := primitive.ObjectIDFromHex(trainingDatumId)
	if err != nil {
		return networks.TrainingDatum{}, errors.Wrap(err, "object id from hex")
	}

	td, err := s.mongo.GetTrainingDatum(ctx, pTrainingDatumId)
	if err != nil {
		return networks.TrainingDatum{}, errors.Wrap(err, "getting training datum")
	}

	return td, nil
}

func (s Service) DeleteTrainingDatum(ctx context.Context, trainingDatumId string) error {
	pTrainingDatumId, err := primitive.ObjectIDFromHex(trainingDatumId)
	if err != nil {
		return errors.Wrap(err, "object id from hex")
	}

	err = s.mongo.DeleteTrainingDatum(ctx, pTrainingDatumId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Wrap(ErrNotFound, err.Error())
		}
		return errors.Wrap(err, "deleting training datum")
	}

	return nil
}
