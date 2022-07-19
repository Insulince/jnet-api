package networks

import (
	"github.com/Insulince/jnet/pkg/network"
	"github.com/Insulince/jnet/pkg/trainer"
)

type (
	TrainingSpec struct {
		NetworkId             string
		Network               network.Network
		Data                  trainer.Data
		TrainingConfiguration trainer.Configuration
	}

	TrainingQueue chan TrainingSpec
)
