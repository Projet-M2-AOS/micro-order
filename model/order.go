package model

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	User         primitive.ObjectID   `json:"user,omitempty" bson:"user,omitempty" validate:"required"`
	Address      string               `json:"address,omitempty" bson:"address,omitempty" validate:"required"`
	Products     []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty" validate:"required"`
	Price        float32              `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	PaymentState string               `json:"paymentState,omitempty" bson:"paymentState,omitempty" validate:"required"`
}

func ValidateStruct(order Order) string {
	strError := ""
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			strError += err.StructNamespace() + " " + err.Tag() + " " + err.Param() + "\n"
		}
	}
	return strError
}
