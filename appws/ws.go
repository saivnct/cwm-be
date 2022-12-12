package appws

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"sol.go/cwm/appfirebase"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/sip"
	"sol.go/cwm/utils"
	"strings"
	"sync"
	"time"
)

type WsCreds struct {
	SessionId string
	User      *model.User
}

//type WsResponse struct {
//	Status  int    `json:"status"`
//	Message string `json:"message"`
//}

type WS struct {
	Server           *socketio.Server
	SendWSMsgChannel chan *WSSendMsg
}

type SocketInfo struct {
	SocketID  string
	SessionID string
}

type WSSendMsg struct {
	Thread *model.SignalThread
	CWMReq *cwmSIPPb.CWMRequest
}

const (
	WS_CRED_AUTH = "authorization"
)

var (
	WSNotFoundSocketErr = errors.New("Not found socket session")
	rWlock              = sync.RWMutex{}

	socketlist   = make(map[string][]SocketInfo) //phoneFull - SocketInfo
	totalSession = 0

	singletonWS *WS
	onceWS      sync.Once
)

func GetWS() *WS {
	onceWS.Do(func() {
		fmt.Println("Init WS...")

		server, err := WSServer()
		if err != nil {
			log.Fatalf("Failed to init WS: %v", err)
		}

		singletonWS = &WS{
			Server: server,
		}
	})
	return singletonWS
}

func (ws *WS) Start() {
	sendWSMsgChannel := make(chan *WSSendMsg, 1000) //make a buffered channel size 1000
	go func(sendWSMsgChannel chan *WSSendMsg) {
		for wsSendMsg := range sendWSMsgChannel {
			err := ws.doSendMsg(wsSendMsg.Thread, wsSendMsg.CWMReq)
			if err != nil {
				log.Println("socketio - sendMsg error", err)
			}
		}
	}(sendWSMsgChannel)

	ws.SendWSMsgChannel = sendWSMsgChannel
}

func (ws *WS) SendMsg(thread *model.SignalThread, req *cwmSIPPb.CWMRequest) {
	ws.SendWSMsgChannel <- &WSSendMsg{
		Thread: thread,
		CWMReq: req,
	}
}

func WSServer() (*socketio.Server, error) {
	wsTransport := websocket.Default
	//pollingTransport := polling.Default

	opts := &engineio.Options{
		PingTimeout:  15 * time.Second,
		PingInterval: 10 * time.Second,
		Transports: []transport.Transport{
			wsTransport,
			//pollingTransport,
		},
	}

	server := socketio.NewServer(opts)

	_, err := server.Adapter(&socketio.RedisAdapterOptions{
		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWD"),
		Prefix:   "socket.io",
	})
	if err != nil {
		log.Fatal("error:", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		//NOTE : Calling Emit inside the OnConnect handler breaks because the connect request is not completed when the handler is called - https://github.com/googollee/go-socket.io/pull/221
		//log.Printf("socketio - OnConnect")

		url := s.URL()
		jwtToken := url.Query().Get(WS_CRED_AUTH)

		jwtPayload, err := hex.DecodeString(jwtToken)
		if err != nil {
			log.Println("WS - Invalid decode jwtToken", err)
			return nil
		}

		claims, _ := utils.ParseJWTToken(string(jwtPayload))
		phoneFull := claims.Subject
		sessionId := claims.ID

		user, _ := dao.GetUserDAO().FindByPhoneFull(context.Background(), phoneFull)

		s.SetContext(WsCreds{
			SessionId: sessionId,
			User:      user,
		})

		s.Join(phoneFull)

		rWlock.Lock() //lock any Writer and Readers

		totalSession++

		sockets, existed := socketlist[phoneFull]
		if !existed {
			sockets = []SocketInfo{}
		}

		sockets = append(sockets, SocketInfo{
			SocketID:  s.ID(),
			SessionID: sessionId,
		})
		socketlist[phoneFull] = sockets

		rWlock.Unlock()

		idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == sessionId })
		if idx >= 0 {
			//update nested obj
			//https://stackoverflow.com/questions/19603542/how-can-i-update-a-property-of-an-object-that-is-contained-in-an-array-of-a-pare
			user.Sessions[idx].Online = true
			updateFields := primitive.M{}
			updateFields["sessions.$[updateSession].online"] = true

			arrayFilter := primitive.M{}
			arrayFilter["updateSession.sessionId"] = sessionId
			arrayFilters := []interface{}{arrayFilter}

			update := primitive.M{"$set": updateFields}
			user, err = dao.GetUserDAO().UpdateByPhoneFull(context.Background(), user.PhoneFull, update, arrayFilters, false)
		}

		log.Printf("socketio - new client of %v connected, socketId: %v, sessionId: %v - total %v session: %v - total system session: %v\n", phoneFull, s.ID(), sessionId, phoneFull, len(sockets), totalSession)

		return nil
	})

	//server.OnEvent("/", "chat2", func(s socketio.Conn, msg string) interface{} {
	//	log.Println("socketio - on chat2:", msg)
	//
	//	wsResponse := WsResponse{
	//		Status:  http.StatusOK,
	//		Message: "OK",
	//	}
	//
	//	server.BroadcastToNamespace("/", "onChatMsg", msg)
	//
	//	return wsResponse
	//})

	server.OnEvent("/", "sendSignalMsg", func(s socketio.Conn, data string) *string {
		//log.Println("socketio - on sendSignalMsg")
		wsCreds, ok := s.Context().(WsCreds)
		if !ok {
			log.Println("socketio - on sendSignalMsg error - can not cast WsCreds")
			return nil
		}

		payload, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Println("socketio - on sendSignalMsg error", err)
			return nil
		}

		cwmRequest := &cwmSIPPb.CWMRequest{}
		err = proto.Unmarshal(payload, cwmRequest)
		if err != nil {
			log.Println("socketio - on sendSignalMsg error", err)
			return nil
		}

		err = sip.VerifySignalMsgFromSipREQ(cwmRequest, wsCreds.User, wsCreds.SessionId)
		if err != nil {
			log.Println("socketio - on sendSignalMsg error", err)
			return nil
		}

		signalThread, cwmRequest, serverDate, err := sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, wsCreds.User)
		if err != nil {
			log.Println("socketio - on sendSignalMsg error", err)
			return nil
		}

		response := cwmSIPPb.CWMResponse{
			Header:  cwmRequest.GetHeader(),
			Code:    200, //TODO - define error code
			Content: utils.Int64ToByteArray(serverDate),
		}

		cwmResponse, err := proto.Marshal(&response)
		if err != nil {
			log.Println("socketio - on sendSignalMsg error", err)
			return nil
		}

		go func(thread *model.SignalThread, req *cwmSIPPb.CWMRequest) {
			GetWS().SendMsg(thread, req)
		}(signalThread, cwmRequest)

		retsult := base64.StdEncoding.EncodeToString(cwmResponse)
		return &retsult
	})

	server.OnEvent("/", "confirmRecieved", func(s socketio.Conn, msgIds string) *string {
		wsCreds, ok := s.Context().(WsCreds)
		if !ok {
			log.Println("socketio - on sendSignalMsg error - can not cast WsCreds")
			return nil
		}

		msgIdList := strings.Split(msgIds, ",")
		//log.Println("socketio - on confirmRecieved", wsCreds.User.PhoneFull, msgIdList)

		go func(msgIdList []string, sessionId string) {
			_, err := sip.ConfirmReceived(msgIdList, sessionId)
			if err != nil {
				log.Println("socketio - confirmRecieved error:", err)
			}
			//log.Println("socketio - confirmRecieved count:", count)

		}(msgIdList, wsCreds.SessionId)

		return &msgIds
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("socketio - OnError:", e)
		//s.Emit("disconnect") // => FORCE CLIENT DISCONNECT
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		//log.Println("socketio - client disconnect:", s.ID(), reason)
		wsCreds, ok := s.Context().(WsCreds)
		if !ok {
			log.Println("socketio - client disconnect can not cast WsCreds")
			return
		}

		rWlock.Lock() //lock any Writer and Readers

		totalSession--

		sockets, existed := socketlist[wsCreds.User.PhoneFull]
		if !existed {
			log.Println("socketio - client disconnect session not found", wsCreds.User.PhoneFull, s.ID())
		} else {
			idx := slices.IndexFunc(sockets, func(socketInfo SocketInfo) bool { return socketInfo.SocketID == s.ID() })
			if idx >= 0 {
				sockets = slices.Delete(sockets, idx, idx+1)
			}

			if len(sockets) == 0 {
				delete(socketlist, wsCreds.User.PhoneFull)
			} else {
				socketlist[wsCreds.User.PhoneFull] = sockets
			}

			log.Printf("socketio - client of %v disconnected, socketId: %v, sessionId: %v - total %v session: %v - total system session: %v\n", wsCreds.User.PhoneFull, s.ID(), wsCreds.SessionId, wsCreds.User.PhoneFull, len(sockets), totalSession)
		}
		rWlock.Unlock()

		user := wsCreds.User

		idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == wsCreds.SessionId })
		if idx >= 0 {
			//update nested obj
			//https://stackoverflow.com/questions/19603542/how-can-i-update-a-property-of-an-object-that-is-contained-in-an-array-of-a-pare
			user.Sessions[idx].Online = false
			updateFields := primitive.M{}
			updateFields["sessions.$[updateSession].online"] = false

			arrayFilter := primitive.M{}
			arrayFilter["updateSession.sessionId"] = wsCreds.SessionId
			arrayFilters := []interface{}{arrayFilter}

			update := primitive.M{"$set": updateFields}
			user, err = dao.GetUserDAO().UpdateByPhoneFull(context.Background(), user.PhoneFull, update, arrayFilters, false)
		}

	})

	return server, nil
}

func (ws *WS) doSendMsg(thread *model.SignalThread, req *cwmSIPPb.CWMRequest) error {
	data, err := proto.Marshal(req)
	if err != nil {
		log.Println("socketio - sendMsg error", err)
		return err
	}

	signalPayload := req.GetContent()
	signalMessage := &cwmSignalMsgPb.SignalMessage{}
	err = proto.Unmarshal(signalPayload, signalMessage)
	if err != nil {
		log.Println("socketio - sendMsg error", err)
		return err
	}

	dataBase64 := base64.StdEncoding.EncodeToString(data)

	shouldPushRemoteWakeup := signalMessage.GetImType() > cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT

	//log.Println("socketio - sendMsg - imtype:", signalMessage.GetImType(), shouldPushRemoteWakeup)

	offlineAndroidSessions := []model.UserSession{}
	offlineIOSSessions := []model.UserSession{}

	rWlock.RLock()
	fromSockets, existed := socketlist[req.GetHeader().GetFrom()] //PhoneFull
	rWlock.RUnlock()
	if existed && len(fromSockets) > 0 {
		//server.BroadcastToNamespace("/", "onChatMsg", data) //client will receive in base64 string
		ws.Server.BroadcastToRoom("/", req.GetHeader().GetFrom(), "onChatMsg", dataBase64) //client will receive in base64 string
	}
	if shouldPushRemoteWakeup {
		fromOfflineAndroidSessions, fromOfflineIosSessions := findOfflineSessions(req.GetHeader().GetFrom(), fromSockets)
		offlineAndroidSessions = append(offlineAndroidSessions, fromOfflineAndroidSessions...)
		offlineIOSSessions = append(offlineIOSSessions, fromOfflineIosSessions...)
	}

	if thread.Type == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO {
		rWlock.RLock()
		toSockets, existed := socketlist[req.GetHeader().GetTo()] //PhoneFull
		rWlock.RUnlock()
		if existed && len(toSockets) > 0 {
			//log.Println("send ws to", req.GetHeader().GetTo())
			ws.Server.BroadcastToRoom("/", req.GetHeader().GetTo(), "onChatMsg", dataBase64) //client will receive in base64 string
		}
		if shouldPushRemoteWakeup {
			toOfflineAndroidSessions, toOfflineIosSessions := findOfflineSessions(req.GetHeader().GetTo(), toSockets)
			offlineAndroidSessions = append(offlineAndroidSessions, toOfflineAndroidSessions...)
			offlineIOSSessions = append(offlineIOSSessions, toOfflineIosSessions...)
		}

	} else if thread.Type == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
		for _, participant := range thread.Participants {
			if participant != req.GetHeader().GetFrom() {
				rWlock.RLock()
				participantSockets, existed := socketlist[participant] //PhoneFull
				rWlock.RUnlock()
				if existed && len(participantSockets) > 0 {
					ws.Server.BroadcastToRoom("/", participant, "onChatMsg", dataBase64) //client will receive in base64 string
				}
				if shouldPushRemoteWakeup {
					participantOfflineAndroidSessions, participantOfflineIosSessions := findOfflineSessions(participant, participantSockets)
					offlineAndroidSessions = append(offlineAndroidSessions, participantOfflineAndroidSessions...)
					offlineIOSSessions = append(offlineIOSSessions, participantOfflineIosSessions...)
				}
			}
		}
	}

	if !shouldPushRemoteWakeup {
		return nil
	}

	////priority low
	if signalMessage.GetImType() > cwmSignalMsgPb.SIGNAL_IM_TYPE_TYPING && signalMessage.GetImType() <= cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT {
		if len(offlineAndroidSessions) > 0 {
			go sendFCMMsg(offlineAndroidSessions, dataBase64, appfirebase.FCM_PRIORITY_NORMAL)
		}
		if len(offlineIOSSessions) > 0 {
			//TODO - send APNS msg
		}
	} else if signalMessage.GetImType() > cwmSignalMsgPb.SIGNAL_IM_TYPE_EVENT {
		if len(offlineAndroidSessions) > 0 {
			go sendFCMMsg(offlineAndroidSessions, dataBase64, appfirebase.FCM_PRIORITY_HIGH)
		}
		if len(offlineIOSSessions) > 0 {
			//TODO - send APNS msg
		}
	}

	return nil
}

func findOfflineSessions(phoneFull string, sockets []SocketInfo) ([]model.UserSession, []model.UserSession) {
	androidSessions := []model.UserSession{}
	iosSessions := []model.UserSession{}

	user, err := dao.GetUserDAO().FindByPhoneFull(context.Background(), phoneFull)
	if err != nil {
		log.Println("checkAndSendPushNotification - error", phoneFull, err)
		return androidSessions, iosSessions
	}

	if sockets == nil || len(sockets) == 0 {
		for _, session := range user.Sessions {
			if session.OsType == grpcCWMPb.OS_TYPE_ANDROID {
				if len(session.PushtokenID) > 0 {
					androidSessions = append(androidSessions, session)
				}
			} else if session.OsType == grpcCWMPb.OS_TYPE_IOS {
				if len(session.PushtokenID) > 0 {
					iosSessions = append(androidSessions, session)
				}
			}
		}
	} else {
		for _, session := range user.Sessions {
			idx := slices.IndexFunc(sockets, func(socketInfo SocketInfo) bool { return socketInfo.SessionID == session.SessionId })
			if idx < 0 {
				if session.OsType == grpcCWMPb.OS_TYPE_ANDROID {
					if len(session.PushtokenID) > 0 {
						androidSessions = append(androidSessions, session)
					}
				} else if session.OsType == grpcCWMPb.OS_TYPE_IOS {
					if len(session.PushtokenID) > 0 {
						iosSessions = append(androidSessions, session)
					}
				}
			}
		}
	}

	return androidSessions, iosSessions
}

func sendFCMMsg(offlineSessions []model.UserSession, dataBase64 string, priorityLevel appfirebase.FCM_PRIORITY) {
	tokens := []string{}
	for _, session := range offlineSessions {
		tokens = append(tokens, session.PushtokenID)
	}

	appfirebase.GetAppFireBase().RequestSendMsg(&appfirebase.FCMSendMsg{
		Tokens:        tokens,
		DataBase64:    dataBase64,
		PriorityLevel: priorityLevel,
	})

}
