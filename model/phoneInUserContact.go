package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhoneInUserContact struct {
	ID        primitive.ObjectID        `json:"_id" bson:"_id,omitempty"`
	PhoneFull string                    `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	Users     map[string]ContactSession `json:"users" bson:"users,omitempty" validate:"required"`
}
