package model

type Contact struct {
	PhoneFull   string   `json:"phoneFull" bson:"pkey,omitempty" validate:"required"`
	ContactName string   `json:"contactName" bson:"contactName,omitempty" validate:"required"`
	Sessions    []string `json:"sessions" bson:"sessions,omitempty" validate:"required"` //user's session that contain this contact
}
