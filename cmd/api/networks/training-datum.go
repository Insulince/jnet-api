package networks

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	TrainingDatum struct {
		Id        primitive.ObjectID `json:"id" bson:"_id"`
		NetworkId primitive.ObjectID `json:"networkId" bson:"networkId"`
		Data      []float64          `json:"data" bson:"data"`
		Truth     []float64          `json:"truth" bson:"truth"`
	}
)
