package networks

import (
	activationfunction "github.com/Insulince/jnet/pkg/activation-function"
	"github.com/Insulince/jnet/pkg/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Network struct {
		Id                 primitive.ObjectID `json:"id" bson:"_id"`
		UpdateKey          string             `json:"updateKey" bson:"updateKey"`
		NeuronMap          []int              `json:"neuronMap" bson:"neuronMap"`
		InputLabels        []string           `json:"inputLabels" bson:"inputLabels"`
		OutputLabels       []string           `json:"outputLabels" bson:"outputLabels"`
		ActivationFunction string             `json:"activationFunction" bson:"activationFunction"`
		Blueprint          []byte             `json:"blueprint" bson:"-"`
		Network            network.Network    `json:"-" bson:"-"`
	}
)

func (nw Network) ToSpec() network.Spec {
	var s network.Spec

	s.NeuronMap = nw.NeuronMap
	s.InputLabels = nw.InputLabels
	s.OutputLabels = nw.OutputLabels
	s.ActivationFunctionName = activationfunction.Name(nw.ActivationFunction)

	return s
}
