package sip

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"log"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/static"
	"sol.go/cwm/utils"
	"time"
)

//var (
//	InvalidSignalEventMsgTypeErr = errors.New("Invalid Signal Event MsgType")
//)

// SOLO MSG
func CreateSignalEventMessageUpdateContactOTT(
	receiver string, receiverSessionsWhiteList []string,
	userOTT *model.User,

) (*cwmSIPPb.CWMRequest, error) {
	cwmRequestHeader := &cwmSIPPb.CWMRequestHeader{
		Method:              cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:                static.ServerEventName,
		To:                  receiver,
		ToSessionsWhiteList: receiverSessionsWhiteList,
	}

	signalEventMessageUpdateContactOTT := &cwmSignalMsgPb.SignalEventMessageUpdateContactOTT{
		PhoneFull:  userOTT.PhoneFull,
		UserId:     userOTT.ID.Hex(),
		Username:   userOTT.Username,
		UserAvatar: userOTT.Avatar,
		FirstName:  userOTT.FirstName,
		LastName:   userOTT.LastName,
	}

	signalEventMessageUpdateContactOTTData, err := proto.Marshal(signalEventMessageUpdateContactOTT)
	if err != nil {
		return nil, err
	}

	signalEventMessage := &cwmSignalMsgPb.SignalEventMessage{
		EventType: cwmSignalMsgPb.SIGNAL_EVENT_MSG_TYPE_UPDATE_CONTACT_OTT,
		Data:      signalEventMessageUpdateContactOTTData,
	}

	signalEventMessageData, err := proto.Marshal(signalEventMessage)
	if err != nil {
		return nil, err
	}
	checksum := fmt.Sprintf("%x", md5.Sum(signalEventMessageData))

	threadId := utils.GetThreadId(static.ServerEventName, receiver)

	now := time.Now()
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT
	signalMessage := &cwmSignalMsgPb.SignalMessage{
		ThreadId:   threadId,
		MsgId:      utils.GenerateUUID(),
		ThreadType: cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO,
		ImType:     imType,
		MsgDate:    now.UnixMilli(),
		ServerDate: now.UnixMilli(),
		Checksum:   checksum,
		Data:       signalEventMessageData,
	}

	signalMessageData, err := proto.Marshal(signalMessage)
	if err != nil {
		return nil, err
	}

	cwmRequest := &cwmSIPPb.CWMRequest{
		Header:  cwmRequestHeader,
		Content: signalMessageData,
	}

	return cwmRequest, nil
}

func CreateSignalEventMessageMsgDelete(sender *model.User, senderSessionID string, signalThread *model.SignalThread, msgIds []string, deleteForAllMembers bool) (*cwmSIPPb.CWMRequest, error) {
	var receiver string
	var threadType = signalThread.Type

	if threadType == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
		receiver = signalThread.ThreadId
	} else {
		idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p != sender.PhoneFull })
		if idx >= 0 {
			receiver = signalThread.Participants[idx]
		} else {
			log.Println("CreateSignalEventMessageMsgDelete - solo thread , but receiver and sender is same person")
			receiver = sender.PhoneFull
		}
	}

	cwmRequestHeader := &cwmSIPPb.CWMRequestHeader{
		Method:        cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:          sender.PhoneFull,
		FromSession:   senderSessionID,
		FromUserName:  sender.Username,
		FromFirstName: sender.FirstName,
		FromLastName:  sender.LastName,
		To:            receiver, //If Group thread -> receiver = ThreadId
	}

	signalEventMessageMsgDelete := &cwmSignalMsgPb.SignalEventMessageMsgDelete{
		ThreadId:            signalThread.ThreadId,
		MsgIds:              msgIds,
		DeleteForAllMembers: deleteForAllMembers,
	}

	signalEventMessageMsgDeleteData, err := proto.Marshal(signalEventMessageMsgDelete)
	if err != nil {
		return nil, err
	}

	signalEventMessage := &cwmSignalMsgPb.SignalEventMessage{
		EventType: cwmSignalMsgPb.SIGNAL_EVENT_MSG_TYPE_MSG_DELETE,
		Data:      signalEventMessageMsgDeleteData,
	}

	signalEventMessageData, err := proto.Marshal(signalEventMessage)
	if err != nil {
		return nil, err
	}
	checksum := fmt.Sprintf("%x", md5.Sum(signalEventMessageData))

	now := time.Now()
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT
	signalMessage := &cwmSignalMsgPb.SignalMessage{
		ThreadId:   signalThread.ThreadId,
		MsgId:      utils.GenerateUUID(),
		ThreadType: threadType,
		ImType:     imType,
		MsgDate:    now.UnixMilli(),
		ServerDate: now.UnixMilli(),
		Checksum:   checksum,
		Data:       signalEventMessageData,
	}

	signalMessageData, err := proto.Marshal(signalMessage)
	if err != nil {
		return nil, err
	}

	cwmRequest := &cwmSIPPb.CWMRequest{
		Header:  cwmRequestHeader,
		Content: signalMessageData,
	}

	return cwmRequest, nil
}

func CreateSignalEventMessageThreadClearMsg(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) (*cwmSIPPb.CWMRequest, error) {
	var receiver string
	var threadType = signalThread.Type

	if threadType == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
		receiver = signalThread.ThreadId
	} else {
		idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p != sender.PhoneFull })
		if idx >= 0 {
			receiver = signalThread.Participants[idx]
		} else {
			log.Println("CreateSignalEventMessageThreadClearMsg - solo thread , but receiver and sender is same person")
			receiver = sender.PhoneFull
		}
	}

	cwmRequestHeader := &cwmSIPPb.CWMRequestHeader{
		Method:        cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:          sender.PhoneFull,
		FromSession:   senderSessionID,
		FromUserName:  sender.Username,
		FromFirstName: sender.FirstName,
		FromLastName:  sender.LastName,
		To:            receiver, //If Group thread -> receiver = ThreadId
	}

	signalEventMessageThreadClearMsg := &cwmSignalMsgPb.SignalEventMessageThreadClearMsg{
		ThreadId:            signalThread.ThreadId,
		DeleteForAllMembers: deleteForAllMembers,
	}

	signalEventMessageThreadClearMsgData, err := proto.Marshal(signalEventMessageThreadClearMsg)
	if err != nil {
		return nil, err
	}

	signalEventMessage := &cwmSignalMsgPb.SignalEventMessage{
		EventType: cwmSignalMsgPb.SIGNAL_EVENT_MSG_TYPE_THREAD_CLEAR_MSG,
		Data:      signalEventMessageThreadClearMsgData,
	}

	signalEventMessageData, err := proto.Marshal(signalEventMessage)
	if err != nil {
		return nil, err
	}
	checksum := fmt.Sprintf("%x", md5.Sum(signalEventMessageData))

	now := time.Now()
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT
	signalMessage := &cwmSignalMsgPb.SignalMessage{
		ThreadId:   signalThread.ThreadId,
		MsgId:      utils.GenerateUUID(),
		ThreadType: threadType,
		ImType:     imType,
		MsgDate:    now.UnixMilli(),
		ServerDate: now.UnixMilli(),
		Checksum:   checksum,
		Data:       signalEventMessageData,
	}

	signalMessageData, err := proto.Marshal(signalMessage)
	if err != nil {
		return nil, err
	}

	cwmRequest := &cwmSIPPb.CWMRequest{
		Header:  cwmRequestHeader,
		Content: signalMessageData,
	}

	return cwmRequest, nil
}

func CreateSignalEventMessageThreadDeleted(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) (*cwmSIPPb.CWMRequest, error) {
	var receiver string
	var threadType = signalThread.Type

	if threadType == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
		receiver = signalThread.ThreadId
	} else {
		idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p != sender.PhoneFull })
		if idx >= 0 {
			receiver = signalThread.Participants[idx]
		} else {
			log.Println("CreateSignalEventMessageThreadDeleted - solo thread , but receiver and sender is same person")
			receiver = sender.PhoneFull
		}
	}

	cwmRequestHeader := &cwmSIPPb.CWMRequestHeader{
		Method:        cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:          sender.PhoneFull,
		FromSession:   senderSessionID,
		FromUserName:  sender.Username,
		FromFirstName: sender.FirstName,
		FromLastName:  sender.LastName,
		To:            receiver, //If Group thread -> receiver = ThreadId
	}

	signalEventMessageThreadDeleted := &cwmSignalMsgPb.SignalEventMessageThreadDeleted{
		ThreadId:            signalThread.ThreadId,
		DeleteForAllMembers: deleteForAllMembers,
	}

	signalEventMessageThreadDeletedData, err := proto.Marshal(signalEventMessageThreadDeleted)
	if err != nil {
		return nil, err
	}

	signalEventMessage := &cwmSignalMsgPb.SignalEventMessage{
		EventType: cwmSignalMsgPb.SIGNAL_EVENT_MSG_TYPE_THREAD_DELETED,
		Data:      signalEventMessageThreadDeletedData,
	}

	signalEventMessageData, err := proto.Marshal(signalEventMessage)
	if err != nil {
		return nil, err
	}
	checksum := fmt.Sprintf("%x", md5.Sum(signalEventMessageData))

	now := time.Now()
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT
	signalMessage := &cwmSignalMsgPb.SignalMessage{
		ThreadId:   signalThread.ThreadId,
		MsgId:      utils.GenerateUUID(),
		ThreadType: threadType,
		ImType:     imType,
		MsgDate:    now.UnixMilli(),
		ServerDate: now.UnixMilli(),
		Checksum:   checksum,
		Data:       signalEventMessageData,
	}

	signalMessageData, err := proto.Marshal(signalMessage)
	if err != nil {
		return nil, err
	}

	cwmRequest := &cwmSIPPb.CWMRequest{
		Header:  cwmRequestHeader,
		Content: signalMessageData,
	}

	return cwmRequest, nil
}

// GROUP MSG
func CreateSignalGroupThreadNotificationMessage(signalThread *model.SignalThread,
	notification_type cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE,
	executor *model.User, executorSessionID string,
	targetMembers []string,
) (*cwmSIPPb.CWMRequest, error) {

	cwmRequestHeader := &cwmSIPPb.CWMRequestHeader{
		Method:        cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:          executor.PhoneFull,
		FromSession:   executorSessionID,
		FromUserName:  executor.Username,
		FromFirstName: executor.FirstName,
		FromLastName:  executor.LastName,
		To:            signalThread.ThreadId,
	}

	signalGroupThreadNotificationMessage := &cwmSignalMsgPb.SignalGroupThreadNotificationMessage{
		NotificationType: notification_type,
		ThreadId:         signalThread.ThreadId,
		Executor:         executor.PhoneFull,
		TargetMembers:    targetMembers,
		GroupName:        signalThread.GroupName,
		Creator:          signalThread.Creator,
		Participants:     signalThread.Participants,
		LastModified:     signalThread.LastModified,
	}

	signalGroupThreadNotificationMessageData, err := proto.Marshal(signalGroupThreadNotificationMessage)
	if err != nil {
		return nil, err
	}

	checksum := fmt.Sprintf("%x", md5.Sum(signalGroupThreadNotificationMessageData))

	now := time.Now()
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_GROUP_THREAD_NOTIFICATION
	protoSignalMessage := &cwmSignalMsgPb.SignalMessage{
		ThreadId:   signalThread.ThreadId,
		MsgId:      utils.GenerateUUID(),
		ThreadType: cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP,
		ImType:     imType,
		MsgDate:    now.UnixMilli(),
		ServerDate: now.UnixMilli(),
		Checksum:   checksum,
		Data:       signalGroupThreadNotificationMessageData,
	}

	protoSignalMessagePayload, err := proto.Marshal(protoSignalMessage)
	if err != nil {
		return nil, err
	}

	cwmRequest := &cwmSIPPb.CWMRequest{
		Header:  cwmRequestHeader,
		Content: protoSignalMessagePayload,
	}

	return cwmRequest, nil
}
