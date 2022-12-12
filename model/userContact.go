package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserContact struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	PhoneFull string             `json:"phoneFull" bson:"pkey,omitempty" validate:"required"` //user
	Contacts  map[string]Contact `json:"contacts" bson:"contacts,omitempty" validate:"required"`
}
