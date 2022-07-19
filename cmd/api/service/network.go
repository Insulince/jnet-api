package service

import (
	"context"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	activationfunction "github.com/Insulince/jnet/pkg/activation-function"
	"github.com/Insulince/jnet/pkg/network"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO(justin): Deactivate a network when it gets update in the DB or deleted (or something else?)

func (s Service) CreateNetwork(ctx context.Context, nm []int, il, ol []string, af string) (string, error) {
	if af == "" {
		af = string(activationfunction.NameSigmoid)
	}

	var spec network.Spec
	spec.NeuronMap = nm
	spec.InputLabels = il
	spec.OutputLabels = ol
	spec.ActivationFunctionName = activationfunction.Name(af)
	jnw, err := network.From(spec)
	if err != nil {
		return "", errors.Wrap(err, "network from spec")
	}

	bp, err := s.networkTranslator.Serialize(jnw)
	if err != nil {
		return "", errors.Wrap(err, "serializing network to bluepring")
	}

	var nw networks.Network
	nw.NeuronMap = nm
	nw.InputLabels = il
	nw.OutputLabels = ol
	nw.ActivationFunction = af
	nw.Blueprint = bp
	id, err := s.mongo.InsertNetwork(ctx, nw)
	if err != nil {
		return "", errors.Wrap(err, "inserting network")
	}

	return id.Hex(), nil
}

func (s Service) SearchNetworks(ctx context.Context) ([]networks.Network, error) {
	nws, err := s.mongo.SearchNetworks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "searching networks")
	}

	return nws, nil
}

func (s Service) GetNetworkById(ctx context.Context, networkId string, checkCache bool) (_ networks.Network, fromCache bool, _ error) {
	if checkCache {
		v, found := s.cache.Get(networkId)
		if found {
			nw, ok := v.(networks.Network)
			if !ok {
				return networks.Network{}, false, errors.Errorf("value in network cache was not of type network, was instead of type %T", v)
			}
			return nw, true, nil
		}
	}

	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return networks.Network{}, false, errors.Wrap(err, "object id from hex")
	}

	nw, err := s.mongo.GetNetworkById(ctx, pNetworkId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return networks.Network{}, false, errors.Wrap(ErrNotFound, err.Error())
		}
		return networks.Network{}, false, errors.Wrap(err, "deserializing network")
	}

	jnw, err := s.networkTranslator.Deserialize(nw.Blueprint)
	if err != nil {
		return networks.Network{}, false, errors.Wrap(err, "deserializing network")
	}
	nw.Network = jnw

	return nw, false, nil
}

// TODO(justin): Need to work on this contract.
func (s Service) UpdateNetwork(ctx context.Context, networkId string, jnw network.Network) error {
	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return errors.Wrap(err, "object id from hex")
	}

	nw, _, err := s.GetNetworkById(ctx, networkId, false)
	if err != nil {
		return errors.Wrap(err, "get network by id")
	}

	bp, err := s.networkTranslator.Serialize(jnw)
	if err != nil {
		return errors.Wrap(err, "serializing network")
	}

	nw.Blueprint = bp

	err = s.mongo.UpsertNetwork(ctx, pNetworkId, nw)
	if err != nil {
		return errors.Wrap(err, "upserting network")
	}

	return nil
}

func (s Service) PatchNetwork(_ context.Context) error {
	return errors.New("not implemented")
}

func (s Service) DeleteNetwork(ctx context.Context, networkId string) error {
	pNetworkId, err := primitive.ObjectIDFromHex(networkId)
	if err != nil {
		return errors.Wrap(err, "object id from hex")
	}

	err = s.mongo.DeleteNetwork(ctx, pNetworkId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Wrap(ErrNotFound, err.Error())
		}
		return errors.Wrap(err, "deleting network")
	}

	err = s.mongo.DeleteAllTrainingData(ctx, pNetworkId)
	if err != nil {
		return errors.Wrap(err, "deleting associated network training data")
	}

	return nil
}

func (s Service) ActivateNetwork(ctx context.Context, networkId string) error {
	nw, _, err := s.GetNetworkById(ctx, networkId, false)
	if err != nil {
		return errors.Wrap(err, "getting network by id")
	}

	jnw, err := s.networkTranslator.Deserialize(nw.Blueprint)
	if err != nil {
		return errors.Wrap(err, "deserializing network")
	}
	nw.Network = jnw

	s.cache.Set(networkId, nw, cache.DefaultExpiration)

	return nil
}

func (s Service) DeactivateNetwork(_ context.Context, networkId string) error {
	s.cache.Delete(networkId)

	return nil
}
