package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sol.go/cwm/proto/cwmSignalMsgPb"
)

type SignalThread struct {
	ID              primitive.ObjectID                `json:"_id" bson:"_id,omitempty"`
	ThreadId        string                            `json:"threadId" bson:"pkey,omitempty" validate:"required"`
	GroupName       string                            `json:"groupName" bson:"groupName,omitempty"`
	Type            cwmSignalMsgPb.SIGNAL_THREAD_TYPE `json:"type" bson:"type,omitempty" validate:"gte=0"`
	AllParticipants []string                          `json:"allParticipants" bson:"allParticipants,omitempty" validate:"required"` //allParticipants included removed members
	Participants    []string                          `json:"participants" bson:"participants,omitempty" validate:"required"`       //current participants of group
	Creator         string                            `json:"creator" bson:"creator,omitempty"`
	Admins          []string                          `json:"admins" bson:"admins,omitempty"`
	CreatedAt       int64                             `json:"createdAt" bson:"createdAt,omitempty" validate:"required"`
	LastModified    int64                             `json:"lastModified" bson:"lastModified,omitempty" validate:"required"`
}
