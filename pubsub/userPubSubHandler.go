package pubsub

import (
	"encoding/base64"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/protobuf/proto"
	"log"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
)

type UserPubSubHandler struct {
}

func (sv *UserPubSubHandler) HandleCWMMsg(data string) {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println("HandleSignalMsg error", err)
		return
	}

	cwmRequest := cwmSIPPb.CWMRequest{}
	err = proto.Unmarshal(payload, &cwmRequest)
	if err != nil {
		log.Println("HandleSignalMsg error", err)
		return
	}
	spew.Dump(&cwmRequest)

	signalPayload := cwmRequest.GetContent()
	signalMessage := cwmSignalMsgPb.SignalMessage{}
	err = proto.Unmarshal(signalPayload, &signalMessage)
	if err != nil {
		log.Println("HandleSignalMsg error", err)
		return
	}
	spew.Dump(&signalMessage)
}
