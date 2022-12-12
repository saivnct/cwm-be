package appgrpc

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/utils"
	"strings"
	"time"
)

func (sv *CWMGRPCService) CreateGroupThread(ctx context.Context, req *grpcCWMPb.CreateGroupThreadRequest) (*grpcCWMPb.CreateGroupThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("CreateGroupThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}
	participants := utils.RemoveSliceDuplicate(req.GetParticipants())

	var threadParticipantInfos []*grpcCWMPb.ThreadParticipantInfo

	includedUser := false
	for _, participant := range participants {
		participantUser, err := dao.GetUserDAO().FindByPhoneFull(ctx, participant)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participants: %v", participant))
		}

		threadParticipantInfo := &grpcCWMPb.ThreadParticipantInfo{
			PhoneFull:  participantUser.PhoneFull,
			UserId:     participantUser.ID.Hex(),
			Username:   participantUser.Username,
			UserAvatar: participantUser.Avatar,
			FirstName:  participantUser.FirstName,
			LastName:   participantUser.LastName,
		}
		threadParticipantInfos = append(threadParticipantInfos, threadParticipantInfo)

		if participant == grpcSession.User.PhoneFull {
			includedUser = true
		}
	}

	if !includedUser {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participants, not include user"))
	}

	now := time.Now().UnixMilli()

	signalThread := &model.SignalThread{
		ThreadId:        utils.GenerateUUID(),
		GroupName:       req.GetGroupName(),
		Type:            cwmSignalMsgPb.SIGNAL_THREAD_TYPE_GROUP,
		AllParticipants: participants,
		Participants:    participants,
		Creator:         grpcSession.User.PhoneFull,
		Admins:          []string{grpcSession.User.PhoneFull},
		CreatedAt:       now,
		LastModified:    now,
	}

	signalThread, err := dao.GetSignalThreadDAO().Save(context.Background(), signalThread)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	for _, participant := range participants {
		_, err := dao.GetUserDAO().UpdateGroupThread(ctx, participant, []string{signalThread.ThreadId}, []string{})
		if err != nil {
			log.Println("CreateGroupThread - UpdateGroupThread failed for", participant, signalThread.ThreadId, err)
		}
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_CREATED,
			executor, executorSessionId,
			[]string{})

		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId)

	return &grpcCWMPb.CreateGroupThreadResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:         signalThread.ThreadId,
			ThreadType:       signalThread.Type,
			GroupName:        signalThread.GroupName,
			Creator:          signalThread.Creator,
			Participants:     signalThread.Participants,
			Admins:           signalThread.Admins,
			LastModified:     signalThread.LastModified,
			ParticipantInfos: threadParticipantInfos,
		},
	}, nil
}

func (sv *CWMGRPCService) CheckGroupThreadInfo(ctx context.Context, req *grpcCWMPb.CheckGroupThreadInfoRequest) (*grpcCWMPb.CheckGroupThreadInfoResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("CheckGroupThreadInfo - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)
	if err != nil {
		log.Println("Invalid thread:", threadId, err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		log.Println("Invalid participant:", threadId)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	participantUsers, err := dao.GetUserDAO().FindByListPhoneFull(ctx, signalThread.Participants)
	if err != nil {
		log.Println("Internal err:", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	var threadParticipantInfos []*grpcCWMPb.ThreadParticipantInfo
	for _, participantUser := range participantUsers {
		threadParticipantInfo := &grpcCWMPb.ThreadParticipantInfo{
			PhoneFull:  participantUser.PhoneFull,
			UserId:     participantUser.ID.Hex(),
			Username:   participantUser.Username,
			UserAvatar: participantUser.Avatar,
			FirstName:  participantUser.FirstName,
			LastName:   participantUser.LastName,
		}
		threadParticipantInfos = append(threadParticipantInfos, threadParticipantInfo)
	}

	return &grpcCWMPb.CheckGroupThreadInfoResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:         signalThread.ThreadId,
			ThreadType:       signalThread.Type,
			GroupName:        signalThread.GroupName,
			Creator:          signalThread.Creator,
			Participants:     signalThread.Participants,
			Admins:           signalThread.Admins,
			LastModified:     signalThread.LastModified,
			ParticipantInfos: threadParticipantInfos,
		},
	}, nil
}

func (sv *CWMGRPCService) ChangeGroupThreadName(ctx context.Context, req *grpcCWMPb.ChangeGroupThreadNameRequest) (*grpcCWMPb.ChangeGroupThreadNameResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("PromoteGroupThreadAdmin - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
	if idxAdmin < 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	}

	updateFields := primitive.M{}
	updateFields["groupName"] = req.GetGroupName()
	updateFields["lastModified"] = time.Now().UnixMilli()
	update := primitive.M{"$set": updateFields}

	signalThread, err = dao.GetSignalThreadDAO().UpdateByThreadId(ctx, threadId, update, []interface{}{}, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_CHANGE_NAME,
			executor, executorSessionId,
			[]string{})
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId)

	return &grpcCWMPb.ChangeGroupThreadNameResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) LeaveGroupThread(ctx context.Context, req *grpcCWMPb.LeaveGroupThreadRequest) (*grpcCWMPb.LeaveGroupThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("LeaveGroupThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	_, err = dao.GetUserDAO().UpdateGroupThread(ctx, grpcSession.User.PhoneFull, []string{}, []string{threadId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateParticipants(ctx, threadId, []string{}, []string{grpcSession.User.PhoneFull})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_LEAVE,
			executor, executorSessionId,
			[]string{})
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId)

	return &grpcCWMPb.LeaveGroupThreadResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) DeleteAndLeaveGroupThread(ctx context.Context, req *grpcCWMPb.DeleteAndLeaveGroupThreadRequest) (*grpcCWMPb.DeleteAndLeaveGroupThreadResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("DeleteAndLeaveGroupThread - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	_, err = dao.GetUserDAO().UpdateGroupThread(ctx, grpcSession.User.PhoneFull, []string{}, []string{threadId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateParticipants(ctx, threadId, []string{}, []string{grpcSession.User.PhoneFull})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	phoneFulls := []string{grpcSession.User.PhoneFull}
	_, err = dao.GetSignalMsgDAO().AppendDeleteUsersByThreadId(ctx, threadId, phoneFulls)
	if err != nil {
		log.Printf("DeleteAndLeaveGroupThread - err: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_LEAVE,
			executor,
			executorSessionId,
			[]string{})
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}

		err = sv.SendEventThreadDeleted(
			executor,
			executorSessionId,
			thread,
			false)

		if err != nil {
			log.Println("Cannot SendEventThreadDeleted", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId)

	return &grpcCWMPb.DeleteAndLeaveGroupThreadResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) AddGroupThreadParticipant(ctx context.Context, req *grpcCWMPb.AddGroupThreadParticipantRequest) (*grpcCWMPb.AddGroupThreadParticipantResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("AddGroupThreadParticipant - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
	if idxAdmin < 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	}

	newParticipants := []string{}
	for _, newParticipant := range req.GetNewParticipants() {
		idxp := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == newParticipant })
		if idxp < 0 {
			_, err = dao.GetUserDAO().UpdateGroupThread(ctx, newParticipant, []string{threadId}, []string{})
			if err == nil {
				newParticipants = append(newParticipants, newParticipant)
			}
		}
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateParticipants(ctx, threadId, newParticipants, []string{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string, targetMembers []string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_ADD_USERS,
			executor, executorSessionId,
			targetMembers)
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId, newParticipants)

	return &grpcCWMPb.AddGroupThreadParticipantResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) RemoveGroupThreadParticipant(ctx context.Context, req *grpcCWMPb.RemoveGroupThreadParticipantRequest) (*grpcCWMPb.RemoveGroupThreadParticipantResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("RemoveGroupThreadParticipant - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
	if idxAdmin < 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	}

	for _, removeParticipant := range req.GetRemoveParticipants() {
		dao.GetUserDAO().UpdateGroupThread(ctx, removeParticipant, []string{}, []string{threadId})
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateParticipants(ctx, threadId, []string{}, req.GetRemoveParticipants())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string, targetMembers []string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_REMOVE_USERS,
			executor, executorSessionId,
			targetMembers)
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId, req.GetRemoveParticipants())

	return &grpcCWMPb.RemoveGroupThreadParticipantResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) PromoteGroupThreadAdmin(ctx context.Context, req *grpcCWMPb.PromoteGroupThreadAdminRequest) (*grpcCWMPb.PromoteGroupThreadAdminResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("PromoteGroupThreadAdmin - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	idx = slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == req.GetPhoneFull() })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid Member: %v", threadId))
	}

	idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
	if idxAdmin < 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	}

	idxAdmin = slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == req.GetPhoneFull() })
	if idxAdmin >= 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Already admin: %v", req.GetPhoneFull()))
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateAdmin(ctx, threadId, []string{req.GetPhoneFull()}, []string{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string, targetMembers []string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_PROMOTE_ADMIN,
			executor, executorSessionId,
			targetMembers)
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId, []string{req.GetPhoneFull()})

	return &grpcCWMPb.PromoteGroupThreadAdminResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}

func (sv *CWMGRPCService) RevokeGroupThreadAdmin(ctx context.Context, req *grpcCWMPb.RevokeGroupThreadAdminRequest) (*grpcCWMPb.RevokeGroupThreadAdminResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("PromoteGroupThreadAdmin - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	if len(strings.TrimSpace(req.PhoneFull)) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid request"))
	}

	threadId := req.GetThreadId()
	signalThread, err := dao.GetSignalThreadDAO().FindByThreadId(ctx, threadId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid thread: %v", threadId))
	}

	idx := slices.IndexFunc(signalThread.Participants, func(p string) bool { return p == grpcSession.User.PhoneFull })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid participant: %v", threadId))
	}

	idx = slices.IndexFunc(signalThread.Participants, func(c string) bool { return c == req.GetPhoneFull() })
	if idx < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid Member: %v", threadId))
	}

	idxAdmin := slices.IndexFunc(signalThread.Admins, func(a string) bool { return a == grpcSession.User.PhoneFull })
	if idxAdmin < 0 {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Dont have thread admin permission: %v", threadId))
	}

	if req.GetPhoneFull() == signalThread.Creator && grpcSession.User.PhoneFull != req.GetPhoneFull() {
		return nil, status.Errorf(codes.PermissionDenied, fmt.Sprintf("Cannot revoke creator: %v", threadId))
	}

	signalThread, err = dao.GetSignalThreadDAO().UpdateAdmin(ctx, threadId, []string{}, []string{req.GetPhoneFull()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(thread *model.SignalThread, executor *model.User, executorSessionId string, targetMembers []string) {
		err := sv.SendGroupThreadNotificationMessage(
			thread,
			cwmSignalMsgPb.SIGNAL_GROUP_THREAD_NOTIFICATION_MSG_TYPE_GROUP_THREAD_REVOKE_ADMIN,
			executor, executorSessionId,
			targetMembers)
		if err != nil {
			log.Println("Cannot SendGroupThreadNotificationMessage", err)
		}
	}(signalThread, grpcSession.User, grpcSession.SessionId, []string{req.GetPhoneFull()})

	return &grpcCWMPb.RevokeGroupThreadAdminResponse{
		GroupThreadInfo: &grpcCWMPb.GroupThreadInfo{
			ThreadId:     signalThread.ThreadId,
			ThreadType:   signalThread.Type,
			GroupName:    signalThread.GroupName,
			Creator:      signalThread.Creator,
			Participants: signalThread.Participants,
			Admins:       signalThread.Admins,
			LastModified: signalThread.LastModified,
		},
	}, nil
}
