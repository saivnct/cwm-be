package appgrpc

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"sol.go/cwm/appws"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/s3Handler"
	"sol.go/cwm/sip"
	"sol.go/cwm/static"
	"sol.go/cwm/utils"
	"strings"
	"time"
)

func (sv *CWMGRPCService) InitialSyncMsg(req *grpcCWMPb.InitialSyncMsgRequest, stream grpcCWMPb.CWMService_InitialSyncMsgServer) error {
	grpcSession, ok := stream.Context().Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("InitialSyncMsg - can not cast GrpcSession")
		return GRPCInvalidSessionErr
	}

	signalThreads, err := dao.GetSignalThreadDAO().FindAllThreadOfUser(context.Background(), grpcSession.User.PhoneFull)
	if err != nil {
		log.Printf("InitialSyncMsg - err: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	log.Printf("InitialSyncMsg - signalThreads: %v\n", len(signalThreads))

	if len(signalThreads) > 0 {
		for _, signalThread := range signalThreads {
			groupThreadInfo := &grpcCWMPb.GroupThreadInfo{
				ThreadId:     signalThread.ThreadId,
				ThreadType:   signalThread.Type,
				GroupName:    signalThread.GroupName,
				Creator:      signalThread.Creator,
				Participants: signalThread.Participants,
				Admins:       signalThread.Admins,
			}
			msgs := []*cwmSIPPb.CWMRequest{}

			signalMsgs, err := dao.GetSignalMsgDAO().FindAllUnreceivedOfThread(
				context.Background(),
				grpcSession.User.PhoneFull,
				grpcSession.SessionId,
				signalThread.ThreadId,
				0, 0,
				0, 10, -1) //Sort descending by createdAt

			if err != nil {
				log.Printf("FetchAllUnreceivedMsgOfThread - err: %v\n", err)

			} else {
				for _, signalMsg := range signalMsgs {
					cwmRequest := &cwmSIPPb.CWMRequest{}
					err = proto.Unmarshal(signalMsg.CwmData, cwmRequest)

					if err != nil {
						continue
					}

					if len(signalMsg.SeenByUsers) > 0 {
						protoSignalPayload := cwmRequest.GetContent()
						protoSignalMessage := &cwmSignalMsgPb.SignalMessage{}
						err = proto.Unmarshal(protoSignalPayload, protoSignalMessage)
						if err != nil {
							continue
						}
						protoSignalMessage.Seenby = signalMsg.SeenByUsers
						updatedProtoSignalPayload, err := proto.Marshal(protoSignalMessage)
						if err != nil {
							continue
						}

						cwmRequest.Content = updatedProtoSignalPayload
					}

					//revert oder of list fetching msg
					//append at 0 index of slice
					msgs = append([]*cwmSIPPb.CWMRequest{cwmRequest}, msgs...)
				}
			}

			res := &grpcCWMPb.InitialSyncMsgResponse{
				GroupThreadInfo: groupThreadInfo,
				Msg:             msgs,
			}
			stream.Send(res)
		}
	}

	return nil
}

func (sv *CWMGRPCService) FetchAllUnreceivedMsg(req *grpcCWMPb.FetchAllUnreceivedMsgRequest, stream grpcCWMPb.CWMService_FetchAllUnreceivedMsgServer) error {
	grpcSession, ok := stream.Context().Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("FetchAllUnreceivedMsg - can not cast GrpcSession")
		return GRPCInvalidSessionErr
	}

	fromDate := req.GetFromDate()
	toDate := req.GetToDate()

	//log.Println("FetchAllUnreceivedMsg ", fromDate, toDate)

	if fromDate <= 0 {
		fromDate = time.Now().AddDate(0, -1, 0).UnixMilli()
	}

	userGroupThreads := grpcSession.User.GroupThreads

	var pageSize int64 = 50

	totalItem, err := dao.GetSignalMsgDAO().CountAllUnreceived(
		context.Background(),
		grpcSession.User.PhoneFull,
		grpcSession.SessionId,
		userGroupThreads,
		fromDate,
		toDate)
	if err != nil {
		log.Printf("FetchAllUnreceivedMsg - err: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	totalPage := totalItem / pageSize
	if totalItem%pageSize > 0 {
		totalPage += 1
	}

	//log.Printf("FetchAllUnreceivedMsg - %v - fromDate %v - toDate %v - totalItem: %v\n", grpcSession.User.PhoneFull, fromDate, toDate, totalItem)

	//log.Printf("FetchAllUnreceivedMsg - totalItem: %v - totalPage: %v - pageSize: %v\n", totalItem, totalPage, pageSize)

	for page := int64(0); page < totalPage; page++ {
		signalMsgs, err := dao.GetSignalMsgDAO().FindAllUnreceived(
			context.Background(),
			grpcSession.User.PhoneFull,
			grpcSession.SessionId,
			userGroupThreads,
			fromDate, toDate,
			page, pageSize)
		if err != nil {
			log.Printf("FetchAllUnreceivedMsg - err: %v\n", err)
			return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}

		for _, signalMsg := range signalMsgs {
			cwmRequest := &cwmSIPPb.CWMRequest{}
			err = proto.Unmarshal(signalMsg.CwmData, cwmRequest)
			if err != nil {
				continue
			}
			//protoSignalMessage := &cwmSignalMsgPb.SignalMessage{}
			//err = proto.Unmarshal(cwmRequest.Content, protoSignalMessage)
			//log.Printf("FetchAllUnreceivedMsg - signalMsg: %v - %v\n", signalMsg.MsgId, string(protoSignalMessage.Data))

			if len(signalMsg.SeenByUsers) > 0 {
				protoSignalPayload := cwmRequest.GetContent()
				protoSignalMessage := &cwmSignalMsgPb.SignalMessage{}
				err = proto.Unmarshal(protoSignalPayload, protoSignalMessage)
				if err != nil {
					continue
				}
				protoSignalMessage.Seenby = signalMsg.SeenByUsers
				updatedProtoSignalPayload, err := proto.Marshal(protoSignalMessage)
				if err != nil {
					continue
				}

				cwmRequest.Content = updatedProtoSignalPayload
			}

			res := &grpcCWMPb.FetchAllUnreceivedMsgResponse{
				Msg: cwmRequest,
			}
			stream.Send(res)
		}
	}

	return nil
}

func (sv *CWMGRPCService) FetchOldMsgOfThread(req *grpcCWMPb.FetchOldMsgOfThreadRequest, stream grpcCWMPb.CWMService_FetchOldMsgOfThreadServer) error {
	grpcSession, ok := stream.Context().Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("FetchOldMsgOfThread - can not cast GrpcSession")
		return GRPCInvalidSessionErr
	}

	threadId := req.GetThreadId()
	if len(threadId) == 0 {
		return status.Errorf(codes.InvalidArgument, "Invalid threadId")
	}
	toDate := req.GetToDate()
	limit := req.GetLimit()
	if limit <= 0 || limit > 10 {
		limit = 10
	}

	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(context.Background(), threadId)
	if err != nil {
		log.Printf("FetchOldMsgOfThread - Invalid thread err: %v\n", err)
		return status.Errorf(codes.InvalidArgument, "Invalid thread")
	}
	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return status.Errorf(codes.PermissionDenied, "Not thread's member")
	}

	//log.Printf("FetchOldMsgOfThread - threadId: %v - toDate: %v - limit: %v\n", threadId, toDate, limit)

	signalMsgs, err := dao.GetSignalMsgDAO().FindOldMsgOfThread(
		context.Background(),
		grpcSession.User.PhoneFull,
		threadId, toDate, limit) //Sort ascending by createdAt
	if err != nil {
		log.Printf("FetchOldMsgOfThread - err: %v\n", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}
	//log.Printf("FetchOldMsgOfThread - items: %v\n", len(signalMsgs))

	//old msgs will be Sorted descending by createdAt, and send to user
	for _, signalMsg := range signalMsgs {
		cwmRequest := &cwmSIPPb.CWMRequest{}
		err = proto.Unmarshal(signalMsg.CwmData, cwmRequest)

		if err != nil {
			log.Printf("FetchOldMsgOfThread - err: %v\n", err)
			break //we must break when error, otherwise logic of client app will be failed
		}

		if len(signalMsg.SeenByUsers) > 0 {
			protoSignalPayload := cwmRequest.GetContent()
			protoSignalMessage := &cwmSignalMsgPb.SignalMessage{}
			err = proto.Unmarshal(protoSignalPayload, protoSignalMessage)
			if err != nil {
				log.Printf("FetchOldMsgOfThread - err: %v\n", err)
				break
			}
			protoSignalMessage.Seenby = signalMsg.SeenByUsers
			updatedProtoSignalPayload, err := proto.Marshal(protoSignalMessage)
			if err != nil {
				log.Printf("FetchOldMsgOfThread - err: %v\n", err)
				break
			}

			cwmRequest.Content = updatedProtoSignalPayload
		}

		res := &grpcCWMPb.FetchOldMsgOfThreadResponse{
			Msg: cwmRequest,
		}
		stream.Send(res)
	}

	return nil
}

func (sv *CWMGRPCService) SendMsg(ctx context.Context, req *grpcCWMPb.SendMsgRequest) (*grpcCWMPb.SendMsgResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("SendMsg - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	cwmRequest := req.GetMsg()

	log.Println("SendMsg")

	err := sip.VerifySignalMsgFromSipREQ(cwmRequest, grpcSession.User, grpcSession.SessionId)
	if err != nil {
		log.Println("SendMsg - invalid msg", err)
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Invalid Msg: %v", err))
	}

	signalThread, cwmRequest, serverDate, err := sip.ParseAndInsertSignalMsgFromSipREQ(cwmRequest, grpcSession.User)
	if err != nil {
		log.Printf("SendMsg - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, req *cwmSIPPb.CWMRequest) {
		appws.GetWS().SendMsg(thread, req)
	}(signalThread, cwmRequest)

	msgResponse := &cwmSIPPb.CWMResponse{
		Header:  cwmRequest.GetHeader(),
		Code:    200, //TODO - define error code
		Content: utils.Int64ToByteArray(serverDate),
	}

	return &grpcCWMPb.SendMsgResponse{
		MsgResponse: msgResponse,
	}, nil
}

func (sv *CWMGRPCService) ConfirmReceivedMsgs(ctx context.Context, req *grpcCWMPb.ConfirmReceivedMsgsRequest) (*grpcCWMPb.ConfirmReceivedMsgsResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("ConfirmReceivedMsgs - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	msgIdList := req.GetMsgIds()
	//log.Printf("ConfirmReceivedMsgs - %v - msgIdList %v", grpcSession.User.PhoneFull, msgIdList)

	_, err := dao.GetSignalMsgDAO().AppendReceivedSessionByMsgIds(context.Background(), msgIdList, grpcSession.SessionId)
	if err != nil {
		log.Printf("ConfirmReceivedMsgs - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	return &grpcCWMPb.ConfirmReceivedMsgsResponse{
		MsgIds: msgIdList,
	}, nil
}

func (sv *CWMGRPCService) DeleteMsgsOfThread(ctx context.Context, req *grpcCWMPb.DeleteMsgsOfThreadRequest) (*grpcCWMPb.DeleteMsgsOfThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("DeleteMsgsOfThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	msgIds := req.GetMsgIds()
	deleteForAllMembers := req.GetDeleteForAllMembers()

	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)
	if err != nil {
		log.Printf("Invalid thread: %v\n", threadId)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	phoneFulls := []string{}

	//if deleteForAllMembers {
	//	if signalThread.Type == grpcCWMPb.SIGNAL_THREAD_TYPE_SIGNAL_THREAD_TYPE_GROUP {
	//		idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.user.PhoneFull })
	//		if idxAdmin < 0 {
	//			return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	//		}
	//	}
	//	phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
	//} else {
	//	phoneFulls = append(phoneFulls, grpcSession.user.PhoneFull)
	//}

	phoneFulls = append(phoneFulls, grpcSession.User.PhoneFull)
	if deleteForAllMembers {
		if signalThread.Type == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
			idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
			if idxAdmin >= 0 {
				phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
			} else {
				log.Printf("DeleteMsgsOfThread: %v\n - want deleteForAllMembers but not admin", threadId)
				deleteForAllMembers = false
			}
		} else {
			phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
		}
	}

	signalMsgs, err := dao.GetSignalMsgDAO().FindAllByMsgIds(ctx, msgIds)
	if err != nil {
		log.Printf("DeleteMsgsOfThread - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	for _, signal := range signalMsgs {
		if signal.ThreadId != threadId {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid msg of thread: %v - %v", signal.MsgId, threadId))
		}
	}

	count, err := dao.GetSignalMsgDAO().AppendDeleteUsersByMsgIds(ctx, msgIds, phoneFulls)
	if err != nil {
		log.Printf("DeleteMsgsOfThread - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(sender *model.User, senderSessionID string, signalThread *model.SignalThread, msgIds []string, deleteForAllMembers bool) {
		err := sv.SendEventMsgsDelete(
			sender,
			senderSessionID,
			signalThread,
			msgIds,
			deleteForAllMembers)
		if err != nil {
			log.Println("Cannot SendEventMsgsDelete", err)
		}
	}(grpcSession.User, grpcSession.SessionId, signalThread, msgIds, deleteForAllMembers)

	return &grpcCWMPb.DeleteMsgsOfThreadResponse{
		Count: count,
	}, nil
}

func (sv *CWMGRPCService) ClearAllMsgOfThread(ctx context.Context, req *grpcCWMPb.ClearAllMsgOfThreadRequest) (*grpcCWMPb.ClearAllMsgOfThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("ClearAllMsgOfThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	deleteForAllMembers := req.GetDeleteForAllMembers()

	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)
	if err != nil {
		log.Printf("Invalid thread: %v\n", threadId)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	phoneFulls := []string{}

	phoneFulls = append(phoneFulls, grpcSession.User.PhoneFull)
	if deleteForAllMembers {
		if signalThread.Type == cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP {
			idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
			if idxAdmin >= 0 {
				phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
			} else {
				log.Printf("ClearAllMsgOfThread: %v\n - want deleteForAllMembers but not admin", threadId)
				deleteForAllMembers = false
			}
		} else {
			phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
		}
	}

	count, err := dao.GetSignalMsgDAO().AppendDeleteUsersByThreadId(ctx, threadId, phoneFulls)
	if err != nil {
		log.Printf("ClearAllMsgOfThread - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) {
		err := sv.SendEventThreadClearMsg(
			sender,
			senderSessionID,
			signalThread,
			deleteForAllMembers)
		if err != nil {
			log.Println("Cannot SendEventThreadClearMsg", err)
		}
	}(grpcSession.User, grpcSession.SessionId, signalThread, deleteForAllMembers)

	return &grpcCWMPb.ClearAllMsgOfThreadResponse{
		Count: count,
	}, nil

}

func (sv *CWMGRPCService) DeleteSoloThread(ctx context.Context, req *grpcCWMPb.DeleteSoloThreadRequest) (*grpcCWMPb.DeleteSoloThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("DeleteSoloThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	deleteForAllMembers := req.GetDeleteForAllMembers()

	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)
	if err != nil {
		log.Printf("Invalid thread: %v\n", threadId)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	if signalThread.Type != cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO {
		log.Printf("Invalid solo thread: %v\n", threadId)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid solo thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	phoneFulls := []string{}
	phoneFulls = append(phoneFulls, grpcSession.User.PhoneFull)
	if deleteForAllMembers {
		phoneFulls = append(phoneFulls, signalThread.AllParticipants...)
	}

	count, err := dao.GetSignalMsgDAO().AppendDeleteUsersByThreadId(ctx, threadId, phoneFulls)
	if err != nil {
		log.Printf("DeleteSoloThread - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(sender *model.User, senderSessionID string, signalThread *model.SignalThread, deleteForAllMembers bool) {
		err := sv.SendEventThreadDeleted(
			sender,
			senderSessionID,
			signalThread,
			deleteForAllMembers)
		if err != nil {
			log.Println("Cannot SendEventThreadDeleted", err)
		}
	}(grpcSession.User, grpcSession.SessionId, signalThread, deleteForAllMembers)

	return &grpcCWMPb.DeleteSoloThreadResponse{
		Count: count,
	}, nil
}

func (sv *CWMGRPCService) UploadMediaMsg(stream grpcCWMPb.CWMService_UploadMediaMsgServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}

	mediaMsgInfo := req.GetMediaMsgInfo()
	fileData := bytes.Buffer{}
	var fileSize int64 = 0

	for {
		//log.Println("waiting to receive more media data")

		req, err = stream.Recv()
		if err == io.EOF {
			log.Println("received file data from client")
			break
		}
		if err != nil {
			log.Println("Cannot receive chunk data", err)
			return err
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		//log.Printf("received a chunk with size: %d\n", size)

		fileSize += int64(size)
		if fileSize > static.MaxFileSize {
			return status.Errorf(codes.InvalidArgument, "file is too large: %d > %d", fileSize, static.MaxFileSize)
		}
		_, err = fileData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	checksum := fmt.Sprintf("%x", md5.Sum(fileData.Bytes()))
	if mediaMsgInfo.Checksum != checksum {
		log.Println("Invalid Checksum", checksum, mediaMsgInfo.Checksum)
		return status.Errorf(codes.InvalidArgument, "Invalid Checksum: %v - %v", checksum, mediaMsgInfo.Checksum)
	}

	mimetype := mimetype.Detect(fileData.Bytes())
	mimetypeString := mimetype.String()
	if mediaMsgInfo.GetMediaType() == cwmSignalMsgPb.SIGNAL_MEDIA_TYPE_IMAGE {
		if !strings.HasPrefix(mimetypeString, "image") {
			return status.Errorf(codes.InvalidArgument, "Invalid media image file")
		}
	} else if mediaMsgInfo.GetMediaType() == cwmSignalMsgPb.SIGNAL_MEDIA_TYPE_VIDEO {
		if !strings.HasPrefix(mimetypeString, "video") {
			return status.Errorf(codes.InvalidArgument, "Invalid video file")
		}
	} else if mediaMsgInfo.GetMediaType() == cwmSignalMsgPb.SIGNAL_MEDIA_TYPE_AUDIO {
		if !strings.HasPrefix(mimetypeString, "audio") {
			return status.Errorf(codes.InvalidArgument, "Invalid audio file")
		}
	} else if mediaMsgInfo.GetMediaType() == cwmSignalMsgPb.SIGNAL_MEDIA_TYPE_DOC {
		if mimetypeString != "application/pdf" &&
			mimetypeString != "application/x-pdf" &&
			mimetypeString != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
			mimetypeString != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" &&
			mimetypeString != "application/vnd.openxmlformats-officedocument.presentationml.presentation" &&
			mimetypeString != "application/msword" &&
			mimetypeString != "application/vnd.ms-word" &&
			mimetypeString != "application/vnd.ms-powerpoint" &&
			mimetypeString != "application/mspowerpoint" &&
			mimetypeString != "application/vnd.ms-excel" &&
			mimetypeString != "application/msexcel" &&
			mimetypeString != "application/vnd.oasis.opendocument.text" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.text" &&
			mimetypeString != "application/vnd.oasis.opendocument.text-template" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.text-template" &&
			mimetypeString != "application/vnd.oasis.opendocument.spreadsheet" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.spreadsheet" &&
			mimetypeString != "application/vnd.oasis.opendocument.spreadsheet-template" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.spreadsheet-template" &&
			mimetypeString != "application/vnd.oasis.opendocument.presentation" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.presentation" &&
			mimetypeString != "application/vnd.oasis.opendocument.presentation-template" &&
			mimetypeString != "application/x-vnd.oasis.opendocument.presentation-template" {
			return status.Errorf(codes.InvalidArgument, "Invalid doc file")
		}
	}

	//log.Println("mimetype", mimetype.String())
	//log.Println("extension", mimetype.Extension())

	fileName := fmt.Sprintf("%s%s%s", static.S3NamePrefxix, checksum, mimetype.Extension())

	//upload to S3
	s3FileStore := s3Handler.GetS3FileStore()
	uploadOutput, err := s3FileStore.UploadFile(fileName, fileData, mimetype)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save file to the s3: %v", err)
	}

	log.Println("Done Upload File to S3... etag:", *uploadOutput.ETag)

	s3FileInfo := &model.S3FileInfo{
		FileId:    utils.GenerateUUID(),
		MsgId:     mediaMsgInfo.MsgId,
		FileName:  fileName,
		FileSize:  fileSize,
		Checksum:  checksum,
		MediaType: mediaMsgInfo.MediaType,
		CreatedAt: time.Now().UnixMilli(),
	}
	s3FileInfo, err = dao.GetS3FileInfoDAO().Save(context.Background(), s3FileInfo)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save s3FileInfo: %v", err)
	}

	//diskFileStore := fileHandler.NewDiskFileStore("media")
	//savedName, err := diskFileStore.SaveToDisk(mediaMsgInfo.FileName, mediaMsgInfo.FileExtension, fileData)
	//if err != nil {
	//	return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	//}

	res := &grpcCWMPb.UploadMediaMsgResponse{
		FileId:   s3FileInfo.FileId,
		FileName: s3FileInfo.FileName,
		FileSize: s3FileInfo.FileSize,
		CheckSum: s3FileInfo.Checksum,
		MsgId:    mediaMsgInfo.MsgId,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	log.Printf("saved file with name: %s, size: %d, MsgId: %s, FileId: %s, checkSum: %s\n", fileName, fileSize, mediaMsgInfo.MsgId, s3FileInfo.FileId, s3FileInfo.Checksum)
	return nil
}

func (sv *CWMGRPCService) DownloadMediaMsg(req *grpcCWMPb.DownloadMediaMsgRequest, stream grpcCWMPb.CWMService_DownloadMediaMsgServer) error {
	_, ok := stream.Context().Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("DownloadMediaMsg - can not cast GrpcSession")
		return GRPCInvalidSessionErr
	}

	msgId := req.GetMsgId()
	fileId := req.GetFileId()

	s3FileInfo, err := dao.GetS3FileInfoDAO().FindByFileId(context.Background(), fileId)
	if s3FileInfo == nil {
		log.Println("Not found s3FileInfo", err)
		return status.Errorf(codes.PermissionDenied, fmt.Sprintf("Not Found s3FileInfo"))
	}

	if s3FileInfo.MsgId != msgId {
		log.Println("Invalid msgId")
		return status.Errorf(codes.PermissionDenied, fmt.Sprintf("Invalid file info"))
	}

	s3Object, err := s3Handler.GetS3FileStore().GetObject(s3FileInfo.FileName)
	if err != nil {
		log.Println("Get s3 Object err", err)
		return status.Errorf(codes.NotFound, fmt.Sprintf("Not found S3 object"))
	}

	if s3Object.Body != nil {
		defer s3Object.Body.Close()
	}
	// If there is no content length, it is a directory
	if s3Object.ContentLength == nil {
		return nil
	}

	//log.Println("DownloadMediaMsg:", *s3Object.ETag, *s3Object.ContentLength)

	buffer := make([]byte, 1024)
	for {
		n, errRead := s3Object.Body.Read(buffer)

		if errRead != nil && errRead != io.EOF {
			return status.Errorf(codes.Internal, fmt.Sprintf("Cannot read S3 object"))
		}

		if n > 0 {
			res := &grpcCWMPb.DownloadMediaMsgResponse{
				ChunkData: buffer[:n],
			}
			err = stream.Send(res)
			if err != nil {
				//To get the real error that contains the gRPC status code, we must call stream.RecvMsg() with a nil parameter. The nil parameter basically means that we don't expect to receive any message, but we just want to get the error that function returns
				log.Println("cannot send chunk to client: ", err)
				return status.Errorf(codes.Internal, fmt.Sprintf("Cannot send chunk to client"))
			}
		}

		if errRead == io.EOF {
			//log.Println("DownloadMediaMsg - reach EOF")
			break
		}
	}

	//log.Println("DownloadMediaMsg - Done")
	return nil
}
