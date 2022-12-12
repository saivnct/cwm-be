package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sol.go/cwm/proto/cwmSignalMsgPb"
)

type SignalMsg struct {
	ID                    primitive.ObjectID                `json:"_id" bson:"_id,omitempty"`
	MsgId                 string                            `json:"msgId" bson:"pkey,omitempty" validate:"required"`
	ThreadId              string                            `json:"threadId" bson:"threadId,omitempty" validate:"required"`
	From                  string                            `json:"from" bson:"from,omitempty" validate:"required"`
	FromSessionsWhiteList []string                          `json:"fromSessionsWhiteList" bson:"fromSessionsWhiteList,omitempty"` //not available in group chat
	To                    string                            `json:"to" bson:"to,omitempty" validate:"required"`                   //not available in group chat
	ToSessionsWhiteList   []string                          `json:"toSessionsWhiteList" bson:"toSessionsWhiteList,omitempty"`     //not available in group chat
	ReceivedSessions      []string                          `json:"receivedSessions" bson:"receivedSessions,omitempty" validate:"required"`
	ThreadType            cwmSignalMsgPb.SIGNAL_THREAD_TYPE `json:"threadType" bson:"threadType,omitempty" validate:"gte=0"`
	DeleteForUsers        []string                          `json:"deleteForUsers" bson:"deleteForUsers,omitempty"`
	SeenByUsers           []string                          `json:"seenByUsers" bson:"seenByUsers,omitempty"`
	CwmData               []byte                            `json:"cwmData" bson:"cwmData,omitempty" validate:"required"` //base64
	CreatedAt             int64                             `json:"createdAt" bson:"createdAt,omitempty" validate:"required"`
}
