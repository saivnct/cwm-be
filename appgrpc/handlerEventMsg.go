package appgrpc

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"log"
	"sol.go/cwm/appws"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/sip"
)

func (sv *CWMGRPCService) SendSMSAuthenCode(phoneFull string) error {
	log.Println("Call SendSMSAuthenCode")

	user, err := dao.GetUserDAO().FindByPhoneFull(context.Background(), phoneFull)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	if user != nil {
		//TODO - find Active Session and send Signal Msg
	} else {
		//TODO - send SMS
	}

	return nil
}

func (sv *CWMGRPCService) SendNotifyUpdateContactOTT(newUserOTT model.User) error {
	log.Println("Call SendNotifyUpdateContactOTT")

	phoneInUserContact, err := dao.GetPhoneInUserContactDAO().FindByPhoneFull(context.Background(), newUserOTT.PhoneFull)
	if err != nil {
		log.Println("SendNotifyUpdateContactOTT err", err)
		return err
	}

	for userPhoneFull, contactSession := range phoneInUserContact.Users {
		user, err := dao.GetUserDAO().FindByPhoneFull(context.Background(), userPhoneFull)
		if err != nil {
			log.Println("SendNotifyUpdateContactOTT err", err)
			continue
		}

		receiverSessions := []string{}

		for _, sessionId := range contactSession.Sessions {
			idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == sessionId })
			if idx < 0 {
				continue
			}

			receiverSessions = append(receiverSessions, sessionId)
		}

		if len(receiverSessions) == 0 {
			log.Println("SendNotifyUpdateContactOTT - No session for", userPhoneFull)
			continue
		}

		cwmRequest, err := sip.CreateSignalEventMessageUpdateContactOTT(
			userPhoneFull, receiverSessions, &newUserOTT)

		if err != nil {
			log.Println("SendNotifyUpdateContactOTT err", err)
			continue
		}

		signalThread, cwmRequest, _, err := sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, nil)
		if err != nil {
			log.Println("SendNotifyUpdateContactOTT err", err)
			continue
		}

		appws.GetWS().SendMsg(signalThread, cwmRequest)

		log.Println("SendNotifyUpdateContactOTT ", userPhoneFull, receiverSessions)
	}

	return nil
}

func (sv *CWMGRPCService) SendGroupThreadNotificationMessage(signalThread *model.SignalThread,
	notification_type cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE,
	executor *model.User, executorSessionID string,
	targetMembers []string,
) error {
	log.Println("Call SendGroupThreadNotificationMessage")
	cwmRequest, err := sip.CreateSignalGroupThreadNotificationMessage(
		signalThread, notification_type,
		executor, executorSessionID,
		targetMembers,
	)

	if err != nil {
		log.Println("SendGroupThreadNotificationMessage err", err)
		return err
	}

	_, cwmRequest, _, err = sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, executor)
	if err != nil {
		log.Println("SendGroupThreadNotificationMessage err", err)
		return err
	}

	appws.GetWS().SendMsg(signalThread, cwmRequest)

	return nil
}

func (sv *CWMGRPCService) SendEventMsgsDelete(sender *model.User, senderSessionID string, signalThread *model.SignalThread, msgIds []string, deleteForAllMembers bool) error {
	log.Println("Call SendEventMsgsDelete")
	cwmRequest, err := sip.CreateSignalEventMessageMsgDelete(sender, senderSessionID, signalThread, msgIds, deleteForAllMembers)

	if err != nil {
		log.Println("SendEventMsgsDelete err", err)
		return err
	}

	_, cwmRequest, _, err = sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, sender)
	if err != nil {
		log.Println("SendEventMsgsDelete err", err)
		return err
	}

	appws.GetWS().SendMsg(signalThread, cwmRequest)

	return nil
}

func (sv *CWMGRPCService) SendEventThreadClearMsg(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) error {
	log.Println("Call SendEventThreadClearMsg")
	cwmRequest, err := sip.CreateSignalEventMessageThreadClearMsg(sender, senderSessionID, signalThread, deleteForAllMembers)

	if err != nil {
		log.Println("SendEventThreadClearMsg err", err)
		return err
	}

	_, cwmRequest, _, err = sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, sender)
	if err != nil {
		log.Println("SendEventThreadClearMsg err", err)
		return err
	}

	appws.GetWS().SendMsg(signalThread, cwmRequest)

	return nil
}

func (sv *CWMGRPCService) SendEventThreadDeleted(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) error {
	log.Println("Call SendEventThreadDeleted")
	cwmRequest, err := sip.CreateSignalEventMessageThreadDeleted(sender, senderSessionID, signalThread, deleteForAllMembers)

	if err != nil {
		log.Println("SendEventThreadDeleted err", err)
		return err
	}

	_, cwmRequest, _, err = sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, sender)
	if err != nil {
		log.Println("SendEventThreadDeleted err", err)
		return err
	}

	appws.GetWS().SendMsg(signalThread, cwmRequest)

	return nil
}
