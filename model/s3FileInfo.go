package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sol.go/cwm/proto/cwmSignalMsgPb"
)

type S3FileInfo struct {
	ID        primitive.ObjectID               `json:"_id" bson:"_id,omitempty"`
	FileId    string                           `json:"fileId" bson:"pkey,omitempty" validate:"required"`
	MsgId     string                           `json:"msgId" bson:"msgId,omitempty" validate:"required"`
	FileName  string                           `json:"fileName" bson:"fileName,omitempty" validate:"required"`
	FileSize  int64                            `json:"fileSize" bson:"fileSize,omitempty" validate:"required"`
	Checksum  string                           `json:"checksum" bson:"checksum,omitempty" validate:"required"`
	MediaType cwmSignalMsgPb.SIGNAL_MEDIA_TYPE `json:"mediaType" bson:"mediaType,omitempty" validate:"gte=0"`
	CreatedAt int64                            `json:"createdAt" bson:"createdAt,omitempty" validate:"required"`
}
