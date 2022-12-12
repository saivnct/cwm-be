package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserTmp struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	PhoneFull        string `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	Phone            string `json:"phone" bson:"phone,omitempty" validate:"required"`
	CountryCode      string `json:"countryCode" bson:"countryCode,omitempty" validate:"required"`
	Authencode       string `json:"authencode" bson:"authencode,omitempty" validate:"required"`
	AuthencodeSendAt int64  `json:"authencodeSendAt" bson:"authencodeSendAt,omitempty" validate:"required"`
	NumberAuthenFail int32  `json:"numberAuthenFail" bson:"numberAuthenFail,omitempty" validate:"gte=0"`
	CreatedAt        int64  `json:"createdAt" bson:"createdAt,omitempty" validate:"required"`
}
