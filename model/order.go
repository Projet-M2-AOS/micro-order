package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	User         primitive.ObjectID   `json:"user,omitempty" bson:"user,omitempty"`
	Address      string               `json:"address,omitempty" bson:"address,omitempty"`
	Products     []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty"`
	Price        float32              `json:"price,omitempty" bson:"price,omitempty"`
	PaymentState string               `json:"paymentState,omitempty" bson:"paymentState,omitempty"`
}
