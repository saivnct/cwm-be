package model

type ContactSession struct {
	PhoneFull string   `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	Sessions  []string `json:"sessions" bson:"sessions,omitempty" validate:"required"`
}
