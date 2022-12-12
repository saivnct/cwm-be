package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhoneLockReg struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	PhoneFull       string `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	Phone           string `json:"phone" bson:"phone,omitempty" validate:"required"`
	CountryCode     string `json:"countryCode" bson:"countryCode,omitempty" validate:"required"`
	NumSmsSend      int32  `json:"numSmsSend" bson:"numSmsSend,omitempty" validate:"gte=0"`
	LastDateSmsSend int64  `json:"lastDateSmsSend" bson:"lastDateSmsSend,omitempty"`
	Locked          bool   `json:"locked" bson:"locked,omitempty"`
}
