package model

import "sol.go/cwm/proto/grpcCWMPb"

type UserSession struct {
	SessionId    string            `json:"sessionId" bson:"sessionId,omitempty"`
	DeviceName   string            `json:"deviceName" bson:"deviceName,omitempty" validate:"required"`
	IMEI         string            `json:"imei" bson:"imei,omitempty" validate:"required"`
	Manufacturer string            `json:"manufacturer" bson:"manufacturer,omitempty" validate:"required"`
	OsType       grpcCWMPb.OS_TYPE `json:"osType" bson:"osType,omitempty" validate:"gte=0"`
	OsVersion    string            `json:"osVersion" bson:"osVersion,omitempty" validate:"required"`
	Online       bool              `json:"online" bson:"online,omitempty"`

	//FCM - APNS_REMOTE
	PushtokenID string `json:"pushtokenID" bson:"pushtokenID,omitempty"`

	//APNS_VOIP
	SecondaryPushtokenID string `json:"secondaryPushtokenID" bson:"secondaryPushtokenID,omitempty"`

	//iOS BundleId
	BundleId string `json:"bundleId" bson:"bundleId,omitempty"`

	//Android AppId
	AppId string `json:"appId" bson:"appId,omitempty"`
}
