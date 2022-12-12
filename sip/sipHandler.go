package sip

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"log"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/utils"
	"time"
)

var (
	InValidSipSignalMsgErr         = errors.New("Invalid SIP Signal Message")
	InValidSipSignalMsgChecksumErr = errors.New("Invalid SIP Signal Message Checksum")
	InValidSipCredentialErr        = errors.New("Invalid SIP Credential")
	InValidThreadErr               = errors.New("Invalid Thread")
)

func VerifySignalMsgFromSipREQ(req *cwmSIPPb.CWMRequest, user *model.User, sessionId string) error {
	if req.GetHeader().GetMethod() != cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE {
		return InValidSipSignalMsgErr
	}

	if req.GetHeader().GetFrom() != user.PhoneFull {
		return InValidSipCredentialErr
	}

	if req.GetHeader().GetFromSession() != sessionId {
		return InValidSipCredentialErr
	}

	signalPayload := req.GetContent()
	signalMessage := &cwmSignalMsgPb.SignalMessage{}
	err := proto.Unmarshal(signalPayload, signalMessage)
	if err != nil {
		return err
	}

	checksum := fmt.Sprintf("%x", md5.Sum(signalMessage.Data))
	if checksum != signalMessage.Checksum {
		return InValidSipSignalMsgChecksumErr
	}

	if signalMessage.ThreadType == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
		signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(context.Background(), signalMessage.GetThreadId())
		if err != nil {
			return err
		}

		idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == user.PhoneFull })
		if idx < 0 {
			log.Println("Invalid participant:", signalThread.ThreadId, user.PhoneFull)
			return InValidSipCredentialErr
		}

	}

	return nil
}

func ParseAndInsertSignalMsgFromSipREQ(req *cwmSIPPb.CWMRequest, fromUser *model.User) (*model.SignalThread, *cwmSIPPb.CWMRequest, int64, error) {

	if req.GetHeader().GetMethod() != cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE {
		return nil, nil, 0, InValidSipSignalMsgErr
	}

	protoSignalPayload := req.GetContent()
	protoSignalMessage := &cwmSignalMsgPb.SignalMessage{}
	err := proto.Unmarshal(protoSignalPayload, protoSignalMessage)
	if err != nil {
		return nil, nil, 0, err
	}

	now := time.Now().UnixMilli()

	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(context.Background(), protoSignalMessage.GetThreadId())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			if protoSignalMessage.GetThreadType() == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO {
				threadId := utils.GetThreadId(req.GetHeader().GetFrom(), req.GetHeader().GetTo())
				if threadId != protoSignalMessage.GetThreadId() {
					log.Println("Invalid threadId", threadId, protoSignalMessage.GetThreadId())
					return nil, nil, 0, InValidThreadErr
				}
				participants := []string{req.GetHeader().GetFrom(), req.GetHeader().GetTo()}

				signalThread = &model.SignalThread{
					ThreadId:        protoSignalMessage.GetThreadId(),
					GroupName:       "",
					Type:            cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO,
					AllParticipants: participants,
					Participants:    participants,
					CreatedAt:       now,
					LastModified:    now,
				}

				signalThread, err = dao.GetSignalThreadDAO().Save(context.Background(), signalThread)
				if err != nil {
					return nil, nil, 0, err
				}
			} else {
				log.Println("Invalid group thread", protoSignalMessage.GetThreadId())
				return nil, nil, 0, InValidThreadErr
			}

		} else {
			return nil, nil, 0, err
		}
	}

	if fromUser != nil {
		req.GetHeader().FromFirstName = fromUser.FirstName
		req.GetHeader().FromLastName = fromUser.LastName
		req.GetHeader().FromUserName = fromUser.Username
	}

	//ignore Typing message -> client will decide what to do with typing msg
	if protoSignalMessage.GetImType() == cwmSignalMsgPb.SIGNAL_IM_TYPE_TYPING {
		return signalThread, req, now, nil
	}

	protoSignalMessage.ServerDate = now
	updatedProtoSignalPayload, err := proto.Marshal(protoSignalMessage)
	if err != nil {
		log.Println("Cannot marshal protoSignalMessage", err)
		return nil, nil, 0, err
	}
	req.Content = updatedProtoSignalPayload

	cwmData, err := proto.Marshal(req)

	receivedSessions := []string{}
	if len(req.GetHeader().GetFromSession()) > 0 {
		receivedSessions = append(receivedSessions, req.GetHeader().GetFromSession())
	}

	signalMsg := &model.SignalMsg{
		MsgId:                 protoSignalMessage.GetMsgId(),
		ThreadId:              protoSignalMessage.GetThreadId(),
		From:                  req.GetHeader().GetFrom(),
		FromSessionsWhiteList: req.GetHeader().GetFromSessionsWhiteList(),
		To:                    req.GetHeader().GetTo(),
		ToSessionsWhiteList:   req.GetHeader().GetToSessionsWhiteList(),
		ReceivedSessions:      receivedSessions,
		ThreadType:            protoSignalMessage.GetThreadType(),
		DeleteForUsers:        []string{},
		SeenByUsers:           []string{},
		CwmData:               cwmData,
		CreatedAt:             now,
	}

	signalMsg, err = dao.GetSignalMsgDAO().Save(context.Background(), signalMsg)
	if err != nil {
		return nil, nil, 0, err
	}

	if protoSignalMessage.GetImType() == cwmSignalMsgPb.SIGNAL_IM_TYPE_SEENSTATE {
		go handleSeenStateMsg(protoSignalMessage, signalMsg.From)
	}

	return signalThread, req, now, nil
}

func ConfirmReceived(msgIdList []string, sessionId string) (int64, error) {
	return dao.GetSignalMsgDAO().AppendReceivedSessionByMsgIds(context.Background(), msgIdList, sessionId)
}

func handleSeenStateMsg(protoSignalMessage *cwmSignalMsgPb.SignalMessage, phoneFull string) {
	//log.Println("handleSeenStateMsg")
	if protoSignalMessage.GetImType() != cwmSignalMsgPb.SIGNAL_IM_TYPE_SEENSTATE {
		return
	}

	protoSignalSeenStateMessage := &cwmSignalMsgPb.SignalSeenStateMessage{}
	err := proto.Unmarshal(protoSignalMessage.GetData(), protoSignalSeenStateMessage)
	if err != nil {
		log.Println(err)
		return
	}

	if protoSignalSeenStateMessage.GetSeenStateType() == cwmSignalMsgPb.SIGNAL_SEENSTATE_MSG_TYPE_SEEN {
		//log.Println("handleSeenStateMsg", protoSignalSeenStateMessage.GetMsgId())
		dao.GetSignalMsgDAO().AppendSeenByUserByMsgIds(context.Background(), protoSignalSeenStateMessage.GetMsgId(), phoneFull)
	}

}
