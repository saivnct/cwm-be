package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sol.go/cwm/static"
	"sol.go/cwm/utils"
	"time"
)

type User struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	PhoneFull     string        `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	Phone         string        `json:"phone" bson:"phone,omitempty" validate:"required"`
	CountryCode   string        `json:"countryCode" bson:"countryCode,omitempty" validate:"required"`
	Username      string        `json:"username" bson:"username,omitempty"`
	Password      string        `json:"password" bson:"password,omitempty" validate:"required"`
	GroupThreads  []string      `json:"groupThreads" bson:"groupThreads,omitempty"` //list PhoneFull + Username
	Sessions      []UserSession `json:"sessions" bson:"sessions,omitempty" validate:"required"`
	BlockThreads  []string      `json:"blockThreads" bson:"blockThreads,omitempty"` //list threadId (Group/SOLO)
	Nonce         string        `json:"nonce" bson:"nonce,omitempty"`
	NonceResponse string        `json:"nonceResponse" bson:"nonceResponse,omitempty"`
	NonceAt       int64         `json:"nonceAt" bson:"nonceAt,omitempty" validate:"required"`

	//info
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
	Avatar    string `json:"avatar" bson:"avatar,omitempty"`
	Thinking  string `json:"thinking" bson:"thinking,omitempty"`
	Gender    int32  `json:"gender" bson:"gender,omitempty"`
	Birthday  int64  `json:"birthday" bson:"birthday,omitempty"`

	CreatedAt int64 `json:"createdAt" bson:"createdAt,omitempty" validate:"required"`
	//UsdtBalance string             `json:"usdtBalance" bson:"usdtBalance,omitempty"`
}

func (user *User) IsExpireNonce() bool {
	now := time.Now()
	nonceAt := time.UnixMilli(user.NonceAt)
	diffMins := now.Sub(nonceAt).Minutes()

	return diffMins > static.NONCETTL
}

func (user *User) IsValidNonce(nonce, response string) bool {
	if user.Nonce != nonce {
		return false
	}

	if user.IsExpireNonce() {
		//log.Printf("%v - Nonce Expire", user.PhoneFull)
		return false
	}

	return user.NonceResponse == response
}

func (user *User) GenNewNonce() {
	nonce := utils.GenerateUUID()
	user.Nonce = nonce
	user.NonceResponse = utils.GetNonceRespone(user.Password, nonce)
	user.NonceAt = time.Now().UnixMilli()

}

func (user *User) ShouldRenewNonce() bool {
	if user.NonceAt == 0 {
		//log.Printf("%v - Should Renew Nonce", user.PhoneFull)
		user.GenNewNonce()
		return true
	}

	now := time.Now()
	nonceAt := time.UnixMilli(user.NonceAt)
	diffMins := now.Sub(nonceAt).Minutes()
	if len(user.Nonce) == 0 || diffMins >= static.TimeToReCreateNONCE {
		//log.Printf("%v - Should Renew Nonce", user.PhoneFull)
		user.GenNewNonce()
		return true
	}

	return false
}
