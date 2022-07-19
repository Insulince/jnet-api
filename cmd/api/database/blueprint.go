package database

import (
	"bytes"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func (m Mongo) Bucket() (*gridfs.Bucket, error) {
	b, err := gridfs.NewBucket(m.JNet())
	if err != nil {
		return nil, errors.Wrap(err, "new bucket")
	}
	return b, nil
}

func (m Mongo) UpsertBlueprint(networkId primitive.ObjectID, bp []byte) error {
	b, err := m.Bucket()
	if err != nil {
		return errors.Wrap(err, "bucket")
	}

	err = b.Delete(networkId)
	if err != nil && !errors.Is(err, gridfs.ErrFileNotFound) {
		return errors.Wrap(err, "deleting file from bucket")
	}

	var buf bytes.Buffer
	buf.Write(bp)
	err = b.UploadFromStreamWithID(networkId, networkId.Hex(), &buf)
	if err != nil {
		return errors.Wrap(err, "upload from stream")
	}

	return nil
}

func (m Mongo) GetBlueprint(networkId primitive.ObjectID) ([]byte, error) {
	b, err := m.Bucket()
	if err != nil {
		return nil, errors.Wrap(err, "bucket")
	}

	buf := bytes.Buffer{}
	_, err = b.DownloadToStream(networkId, &buf)
	if err != nil {
		return nil, errors.Wrap(err, "download to stream")
	}
	bp := buf.Bytes()

	return bp, nil
}

func (m Mongo) DeleteBlueprint(networkId primitive.ObjectID) error {
	b, err := m.Bucket()
	if err != nil {
		return errors.Wrap(err, "bucket")
	}

	err = b.Delete(networkId)
	if err != nil {
		return errors.Wrap(err, "deleting file from bucket")
	}

	return nil
}
