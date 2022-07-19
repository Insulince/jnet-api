package database

import (
	"context"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m Mongo) TrainingData() *mongo.Collection {
	collection := m.JNet().Collection("training-data")

	return collection
}

func (m Mongo) InsertTrainingData(ctx context.Context, tds []networks.TrainingDatum) ([]primitive.ObjectID, error) {
	var documents []interface{}
	for _, td := range tds {
		td.Id = primitive.NewObjectID()
		documents = append(documents, td)
	}

	r, err := m.TrainingData().InsertMany(ctx, documents)
	if err != nil {
		return nil, errors.Wrap(err, "inserting many")
	}

	var ids []primitive.ObjectID
	for _, insertedID := range r.InsertedIDs {
		id, ok := insertedID.(primitive.ObjectID)
		if !ok {
			return nil, errors.Errorf("cannot cast resulting id from inserted document to a primitive.ObjectID, was instead of type %T", insertedID)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (m Mongo) SearchTrainingData(ctx context.Context, networkId primitive.ObjectID) ([]networks.TrainingDatum, error) {
	cur, err := m.TrainingData().Find(ctx, bson.D{{"networkId", networkId}})
	if err != nil {
		return nil, errors.Wrap(err, "finding documents")
	}
	defer func() { _ = cur.Close(ctx) }()

	var tds []networks.TrainingDatum
	for cur.Next(ctx) {
		var td networks.TrainingDatum
		if err := cur.Decode(&td); err != nil {
			return nil, errors.Wrap(err, "decoding document")
		}
		tds = append(tds, td)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor")
	}

	return tds, nil
}

func (m Mongo) DeleteAllTrainingData(ctx context.Context, networkId primitive.ObjectID) error {
	_, err := m.TrainingData().DeleteMany(ctx, bson.D{{"networkId", networkId}})
	if err != nil {
		return errors.Wrap(err, "deleting many")
	}

	return nil
}

func (m Mongo) GetTrainingDatum(ctx context.Context, trainingDatumId primitive.ObjectID) (networks.TrainingDatum, error) {
	var td networks.TrainingDatum
	err := m.TrainingData().FindOne(ctx, bson.D{{"_id", trainingDatumId}}).Decode(&td)
	if err != nil {
		return networks.TrainingDatum{}, errors.Wrap(err, "decoding document")
	}

	return td, nil
}

func (m Mongo) DeleteTrainingDatum(ctx context.Context, trainingDatumId primitive.ObjectID) error {
	_, err := m.TrainingData().DeleteOne(ctx, bson.D{{"_id", trainingDatumId}})
	if err != nil {
		return errors.Wrap(err, "deleting one")
	}

	return nil
}
