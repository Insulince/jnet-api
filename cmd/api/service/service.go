package service

import (
	"github.com/Insulince/jnet-api/untitled/cmd/api/database"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet/pkg/network"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type (
	Service struct {
		mongo             database.Mongo
		networkTranslator network.Translator
		trainingQueue     networks.TrainingQueue
		cache             *cache.Cache
	}
)

func New(m database.Mongo, nt network.Translator, tq networks.TrainingQueue) Service {
	var s Service

	c := cache.New(1*time.Hour, 1*time.Minute)

	s.mongo = m
	s.networkTranslator = nt
	s.trainingQueue = tq
	s.cache = c

	return s
}
