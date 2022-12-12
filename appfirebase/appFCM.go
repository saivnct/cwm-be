package appfirebase

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"log"
	"sync"
)

type FCM_PRIORITY int32

const (
	FCM_PRIORITY_NORMAL FCM_PRIORITY = 0
	FCM_PRIORITY_HIGH   FCM_PRIORITY = 1

	DataMsgKey = "msg"
)

type AppFireBase struct {
	App               firebase.App
	SendFCMMsgChannel chan *FCMSendMsg
}

type FCMSendMsg struct {
	Tokens        []string
	DataBase64    string
	PriorityLevel FCM_PRIORITY
}

var (
	InvalidOSTypeErr = errors.New("Invalid OS TYPE")
	singletonAppFCM  *AppFireBase
	onceAppFCM       sync.Once
)

func GetAppFireBase() *AppFireBase {
	onceAppFCM.Do(func() {
		fmt.Println("Init App Firebase...")
		app, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalf("error initializing app firebase: %v\n", err)
		}

		singletonAppFCM = &AppFireBase{
			App: *app,
		}
	})
	return singletonAppFCM
}

func (appFirebase *AppFireBase) Start() {
	sendFCMMsgChannel := make(chan *FCMSendMsg, 1000) //make a buffered channel size 1000
	go func(sendFCMMsgChannel chan *FCMSendMsg) {
		for fcmSendMsg := range sendFCMMsgChannel {
			appFirebase.handleSendMsg(fcmSendMsg)
		}
	}(sendFCMMsgChannel)

	appFirebase.SendFCMMsgChannel = sendFCMMsgChannel
}

func (appFirebase *AppFireBase) RequestSendMsg(fcmSendMsg *FCMSendMsg) {
	appFirebase.SendFCMMsgChannel <- fcmSendMsg
}

func (appFirebase *AppFireBase) handleSendMsg(fcmSendMsg *FCMSendMsg) error {
	if len(fcmSendMsg.Tokens) == 1 {
		tokenID := fcmSendMsg.Tokens[0]
		err := appFirebase.sendSingleMsg(tokenID, fcmSendMsg.DataBase64, fcmSendMsg.PriorityLevel)
		if err != nil {
			return err
		}

	} else {
		tokens := []string{}
		for _, tokenID := range fcmSendMsg.Tokens {
			tokens = append(tokens, tokenID)

			if len(tokens) == 500 { //FCM APIs allow up to 500 device registration tokens per invocation
				appFirebase.sendBatchMsg(tokens, fcmSendMsg.DataBase64, fcmSendMsg.PriorityLevel)
				tokens = []string{}
			}
		}

		if len(tokens) > 0 {
			appFirebase.sendBatchMsg(tokens, fcmSendMsg.DataBase64, fcmSendMsg.PriorityLevel)
		}

	}

	return nil
}

func (appFirebase *AppFireBase) sendSingleMsg(token string, dataBase64 string, priorityLevel FCM_PRIORITY) error {
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := appFirebase.App.Messaging(ctx)
	if err != nil {
		return err
	}

	priority := "normal"
	if priorityLevel == FCM_PRIORITY_HIGH {
		priority = "high"
	}

	androidConfig := &messaging.AndroidConfig{
		//CollapseKey:           "",
		Priority: priority, // one of "normal" or "high"
	}

	message := &messaging.Message{
		Data: map[string]string{
			DataMsgKey: dataBase64,
		},
		Android: androidConfig,
		Token:   token,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	// Response is a message ID string.

	_, err = client.Send(ctx, message)

	if err != nil {
		log.Println("sendFCMMsg failed", err)
	} else {
		log.Printf("sendFCMMsg success to %v - priorityLevel: %v", token, priorityLevel)
	}

	return err
}

func (appFirebase *AppFireBase) sendBatchMsg(tokens []string, dataBase64 string, priorityLevel FCM_PRIORITY) (*messaging.BatchResponse, error) {
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := appFirebase.App.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	priority := "normal"
	if priorityLevel == FCM_PRIORITY_HIGH {
		priority = "high"
	}

	androidConfig := &messaging.AndroidConfig{
		//CollapseKey:           "",
		Priority: priority, // one of "normal" or "high"
	}

	// See documentation on defining a message payload.
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			DataMsgKey: dataBase64,
		},
		Android: androidConfig,
		Tokens:  tokens,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	// Response is a message ID string.
	batchResponse, err := client.SendMulticast(ctx, message)

	if err != nil {
		log.Println("sendFCMMsg failed", err)
	} else {
		log.Printf("SendBatchMsg - success %v - failed %v - of all %v - priorityLevel: %v", batchResponse.SuccessCount, batchResponse.FailureCount, len(tokens), priorityLevel)
	}

	return batchResponse, err
}
