package database

import (
	"context"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m Mongo) Networks() *mongo.Collection {
	collection := m.JNet().Collection("networks")

	return collection
}

func (m Mongo) InsertNetwork(ctx context.Context, nw networks.Network) (primitive.ObjectID, error) {
	nw.Id = primitive.NewObjectID()

	updateKey, err := randomHex(updateKeyLength)
	if err != nil {
		return primitive.ObjectID{}, errors.Wrap(err, "random hex")
	}
	nw.UpdateKey = updateKey

	err = m.UpsertBlueprint(nw.Id, nw.Blueprint)
	if err != nil {
		return primitive.ObjectID{}, errors.Wrap(err, "inserting blueprint")
	}

	r, err := m.Networks().InsertOne(ctx, nw)
	if err != nil {
		return primitive.ObjectID{}, errors.Wrap(err, "inserting one document")
	}

	id, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, errors.Errorf("cannot cast resulting id from inserted document to a primitive.ObjectID, was instead of type %T", r.InsertedID)
	}

	return id, nil
}

// TODO(justin): Hide blueprints option
func (m Mongo) SearchNetworks(ctx context.Context) ([]networks.Network, error) {
	cur, err := m.Networks().Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "finding documents")
	}
	defer func() { _ = cur.Close(ctx) }()

	var nws []networks.Network
	for cur.Next(ctx) {
		var nw networks.Network
		if err := cur.Decode(&nw); err != nil {
			return nil, errors.Wrap(err, "decoding document")
		}

		bp, err := m.GetBlueprint(nw.Id)
		if err != nil {
			return nil, errors.Wrap(err, "getting blueprint")
		}
		nw.Blueprint = bp

		nws = append(nws, nw)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor")
	}

	return nws, nil
}

func (m Mongo) GetNetworkById(ctx context.Context, networkId primitive.ObjectID) (networks.Network, error) {
	var nw networks.Network
	err := m.Networks().FindOne(ctx, bson.D{{"_id", networkId}}).Decode(&nw)
	if err != nil {
		return networks.Network{}, errors.Wrap(err, "finding network by id")
	}

	bp, err := m.GetBlueprint(nw.Id)
	if err != nil {
		return networks.Network{}, errors.Wrap(err, "getting blueprint")
	}
	nw.Blueprint = bp

	return nw, nil
}

func (m Mongo) UpsertNetwork(ctx context.Context, networkId primitive.ObjectID, nw networks.Network) error {
	nw.Id = networkId

	err := m.UpsertBlueprint(nw.Id, nw.Blueprint)
	if err != nil {
		return errors.Wrap(err, "upserting blueprint")
	}

	_, err = m.Networks().UpdateByID(ctx, networkId, bson.M{"$set": nw})
	if err != nil {
		return errors.Wrap(err, "update by id")
	}

	return nil
}

func (m Mongo) DeleteNetwork(ctx context.Context, networkId primitive.ObjectID) error {
	err := m.DeleteBlueprint(networkId)
	if err != nil {
		return errors.Wrap(err, "deleting blueprint")
	}

	_, err = m.Networks().DeleteOne(ctx, bson.D{{"_id", networkId}})
	if err != nil {
		return errors.Wrap(err, "deleting one")
	}

	return nil
}
