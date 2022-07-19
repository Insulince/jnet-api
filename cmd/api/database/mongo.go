package database

import (
	"context"
	"encoding/hex"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
)

const (
	updateKeyLength = 24
)

type (
	Mongo struct {
		client *mongo.Client
	}
)

func NewMongo(ctx context.Context, client *mongo.Client) (Mongo, error) {
	var m Mongo

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return Mongo{}, errors.Wrap(err, "pinging mongo server")
	}

	m.client = client

	return m, nil
}

func (m Mongo) JNet() *mongo.Database {
	database := m.client.Database("jnet")

	return database
}

func randomHex(n int) (string, error) {
	// divide by two because there are two hexadecimal characters per byte, but n should be the number of characters.
	bytes := make([]byte, n/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.Wrap(err, "reading random bytes")
	}
	return hex.EncodeToString(bytes), nil
}
