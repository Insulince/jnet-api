package config

import (
	"github.com/pkg/errors"
	"os"
)

const (
	EnvMongoConnectionString = "MONGO_CONNECTION_STRING"
)

type (
	Config struct {
		MongoConnectionString string
	}
)

var (
	ErrBadEnv = errors.Errorf("bad env")
)

func GetConfig() (Config, error) {
	var c Config

	mcs, found := os.LookupEnv(EnvMongoConnectionString)
	if !found {
		return Config{}, badEnv(EnvMongoConnectionString)
	}
	c.MongoConnectionString = mcs

	return c, nil
}

func badEnv(envVar string) error {
	return errors.Wrapf(ErrBadEnv, "missing %s", envVar)
}
