package appgrpc

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/static"
	"sol.go/cwm/utils"
	"strings"
	"time"
)

func (sv *CWMGRPCService) CreatUser(ctx context.Context, req *grpcCWMPb.CreatAccountRequest) (*grpcCWMPb.CreatAccountResponse, error) {
	//TODO - check phonenumber pattern & country code

	//userCtx, ok := ctx.Value(GRPC_CTX_KEY_USER).(*model.User)
	//if ok {
	//	log.Println("receive request from", userCtx.PhoneFull)
	//}

	if len(req.GetPhone()) == 0 || len(req.GetCountryCode()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Phone Number")
	}

	phoneFull := utils.GetPhoneFull(req.GetCountryCode(), req.GetPhone())

	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.Canceled, "the client canceled the request")
	}

	phoneLockReg, err := dao.GetPhoneLockRegDAO().FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //phoneLockReg == nil
			phoneLockReg = &model.PhoneLockReg{
				PhoneFull:       phoneFull,
				Phone:           req.GetPhone(),
				CountryCode:     req.GetCountryCode(),
				NumSmsSend:      0,
				LastDateSmsSend: 0,
				Locked:          false,
			}
		} else {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}

	authenCode := utils.GenerateAuthenCode()
	log.Println("authenCode", authenCode)
	now := time.Now()

	userTmp, err := dao.GetUserTmpDAO().FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //userTmp == nil
			userTmp = &model.UserTmp{
				PhoneFull:   phoneFull,
				Phone:       req.GetPhone(),
				CountryCode: req.GetCountryCode(),
				CreatedAt:   now.UnixMilli(),
			}
		} else {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}
	userTmp.Authencode = authenCode
	userTmp.AuthencodeSendAt = now.UnixMilli()
	userTmp.NumberAuthenFail = 0

	if phoneLockReg.NumSmsSend > static.MaxResendAuthenCode {
		phoneLockReg.Locked = true
	}

	if phoneLockReg.Locked {
		lastDateSmsSend := time.UnixMilli(phoneLockReg.LastDateSmsSend)
		diffMins := now.Sub(lastDateSmsSend).Minutes()
		//fmt.Println("diffMins", diffMins)
		if diffMins >= static.TimeResetPhoneCreateAccountLock {
			phoneLockReg.Locked = false
			phoneLockReg.NumSmsSend = 0
		} else {
			log.Println("Over Max Resend AuthenCode", phoneFull)
			return nil, status.Errorf(codes.Aborted, "Over Max Resend AuthenCode")
		}
	}
	phoneLockReg.NumSmsSend++
	phoneLockReg.LastDateSmsSend = now.UnixMilli()

	userTmp, err = dao.GetUserTmpDAO().Save(ctx, userTmp)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal err: %v", err),
		)
	}

	phoneLockReg, err = dao.GetPhoneLockRegDAO().Save(ctx, phoneLockReg)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal err: %v", err),
		)
	}

	go func(phoneFull string) {
		err := sv.SendSMSAuthenCode(phoneFull)
		if err != nil {
			log.Println("Cannot send SMS", err)
		}
	}(userTmp.PhoneFull)

	creatAccountResponse := grpcCWMPb.CreatAccountResponse{
		PhoneFull:         phoneFull,
		AuthenCodeTimeOut: static.AuthencodeTimeOut,
		NumberPrefix:      static.NumberPrefix,
	}

	return &creatAccountResponse, nil
}

func (sv *CWMGRPCService) VerifyAuthencode(ctx context.Context, req *grpcCWMPb.VerifyAuthencodeRequest) (*grpcCWMPb.VerifyAuthencodeResponse, error) {
	if len(req.GetPhone()) == 0 || len(req.GetCountryCode()) == 0 || len(req.GetAuthencode()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Arguments")
	}

	if req.GetDeviceInfo() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid DeviceInfo")
	}

	userSession := model.UserSession{
		DeviceName:           req.GetDeviceInfo().GetDeviceName(),
		IMEI:                 req.GetDeviceInfo().GetImei(),
		Manufacturer:         req.GetDeviceInfo().GetManufacturer(),
		OsType:               req.GetDeviceInfo().GetOs(),
		OsVersion:            req.GetDeviceInfo().GetOsVersion(),
		PushtokenID:          "",
		SecondaryPushtokenID: "",
		BundleId:             "",
		AppId:                "",
	}
	userSession.SessionId = utils.GenerateUUID()

	//spew.Dump(userSession)

	validate := validator.New()
	err := validate.Struct(userSession)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid DeviceInfo: %v", err),
		)
	}

	phoneFull := utils.GetPhoneFull(req.GetCountryCode(), req.GetPhone())
	userTmp, _ := dao.GetUserTmpDAO().FindByPhoneFull(ctx, phoneFull)
	if userTmp == nil {
		return nil, status.Errorf(codes.Aborted, "Not found UserTmp")
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.Canceled, "the client canceled the request")
	}

	now := time.Now()
	authencodeSendAt := time.UnixMilli(userTmp.AuthencodeSendAt)
	if now.Sub(authencodeSendAt).Seconds() > static.AuthencodeTimeOut {
		return nil, status.Errorf(codes.Aborted, "Authencode TimeOut")
	}

	if userTmp.Authencode != req.GetAuthencode() {
		userTmp.NumberAuthenFail++
		if userTmp.NumberAuthenFail > static.MaxNumberAuthencodeFail {
			_, err = dao.GetUserTmpDAO().DeleteByPhoneFull(ctx, phoneFull)
			if err != nil {
				log.Println(err)
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
			}
			return nil, status.Error(codes.Aborted, "Max Number Authencode Fail")
		}
		_, err = dao.GetUserTmpDAO().Save(ctx, userTmp)
		if err != nil {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
		return nil, status.Error(codes.Aborted, "Invalid Authencode")
	}

	user, err := dao.GetUserDAO().FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //user == nil
			password := utils.GenerateUUID()
			user = &model.User{
				PhoneFull:    phoneFull,
				Phone:        req.GetPhone(),
				CountryCode:  req.GetCountryCode(),
				Password:     password,
				Sessions:     []model.UserSession{},
				GroupThreads: []string{},
				BlockThreads: []string{},
				CreatedAt:    now.UnixMilli(),
			}
		} else {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}

	authencodeSHA256 := sha256.Sum256([]byte(userTmp.Authencode))
	passEnc, iv, err := utils.Aes256Encode([]byte(user.Password), authencodeSHA256[:])
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	_, err = dao.GetUserTmpDAO().DeleteByPhoneFull(ctx, phoneFull)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.IMEI == userSession.IMEI })
	if idx >= 0 {
		log.Println("IMEI existed", userSession.IMEI)
		//userSession.SessionId = user.Sessions[idx].SessionId	//keep sessionId for same IMEI
		user.Sessions[idx] = userSession
	} else {
		log.Println("IMEI not existed", userSession.IMEI)
		user.Sessions = append(user.Sessions, userSession)
	}

	user.ShouldRenewNonce()

	user, err = dao.GetUserDAO().Save(ctx, user)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(newUserOTT model.User) {
		err := sv.SendNotifyUpdateContactOTT(newUserOTT)
		if err != nil {
			log.Println("Cannot SendNotifyUpdateContactOTT", err)
		}
	}(*user)

	deviceinfo := req.GetDeviceInfo()
	deviceinfo.SessionId = &userSession.SessionId
	verifyAuthencodeResponse := grpcCWMPb.VerifyAuthencodeResponse{
		PassEnc:    passEnc,
		Iv:         iv,
		DeviceInfo: deviceinfo,
		Nonce:      user.Nonce,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.Username,
	}

	return &verifyAuthencodeResponse, nil
}

func (sv *CWMGRPCService) Login(ctx context.Context, req *grpcCWMPb.LoginRequest) (*grpcCWMPb.LoginResponse, error) {

	if len(req.GetPhoneFull()) == 0 || len(req.GetSessionId()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Arguments")
	}

	log.Println("Call Login", req.GetPhoneFull())

	user, err := dao.GetUserDAO().FindByPhoneFull(context.Background(), req.GetPhoneFull())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //NotFound user
			return nil, status.Errorf(codes.NotFound, "NotFound user")
		} else {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}

	idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == req.GetSessionId() })
	if idx < 0 {
		return nil, status.Errorf(codes.NotFound, "NotFound sessionId")
	}

	isValidNonce := user.IsValidNonce(req.GetNonce(), req.GetResponse())
	if !isValidNonce {
		log.Printf("GRPC login - %v inValidNonce!!!!!", user.PhoneFull, user.NonceResponse, req.GetResponse())
		shouldRenewNonce := user.ShouldRenewNonce()
		if shouldRenewNonce {
			log.Printf("GRPC login - %v shouldRenewNonce!!!!!", user.PhoneFull)
			user, err = dao.GetUserDAO().UpdateNonce(ctx, user.PhoneFull)
			if err != nil {
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
			}
		}

		//NOT SUPPORT FOR JAVA PROTOBUF_LITE
		//fieldViolation := &errdetails.BadRequest_FieldViolation{
		//	Field:       "nonce",
		//	Description: user.Nonce,
		//}
		//badRequest := &errdetails.BadRequest{}
		//badRequest.FieldViolations = append(badRequest.FieldViolations, fieldViolation)
		//
		//st := status.New(codes.Unauthenticated, "invalid username")
		//st, err = st.WithDetails(badRequest)
		//if err != nil {
		//	return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		//}
		//
		//return nil, st.Err()
		return nil, status.Errorf(codes.Unauthenticated, user.Nonce)
	}

	expirationTime := jwt.NumericDate{Time: time.Now().Add(static.JWTTTL * time.Minute)}
	// Create the JWT claims, which includes the username and expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &expirationTime,
		Subject:   user.PhoneFull,
		ID:        req.GetSessionId(),
	})

	responseJWT, err := token.SignedString(static.JWTKey())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	loginResponse := &grpcCWMPb.LoginResponse{
		Jwt:    responseJWT,
		JwtTTL: static.JWTTTL,
	}

	return loginResponse, nil
}

func (sv *CWMGRPCService) UpdateProfile(ctx context.Context, req *grpcCWMPb.UpdateProfileRequest) (*grpcCWMPb.UpdateProfileResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("UpdateProfile - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	//log.Println("UpdateProfile", req.GetFirstName(), req.GetLastName())

	firstName := strings.TrimSpace(req.GetFirstName())
	lastName := strings.TrimSpace(req.GetLastName())
	//log.Println("UpdateProfile", firstName, lastName)

	updateFields := primitive.M{}
	if len(firstName) > 0 {
		updateFields["firstName"] = firstName
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid firstName")
	}

	if len(lastName) > 0 {
		updateFields["lastName"] = lastName
	}
	update := primitive.M{"$set": updateFields}

	user, err := dao.GetUserDAO().UpdateByPhoneFull(ctx, grpcSession.User.PhoneFull, update, []interface{}{}, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(userOTT model.User) {
		err := sv.SendNotifyUpdateContactOTT(userOTT)
		if err != nil {
			log.Println("Cannot SendNotifyUpdateContactOTT", err)
		}
	}(*user)

	return &grpcCWMPb.UpdateProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (sv *CWMGRPCService) UpdateUsername(ctx context.Context, req *grpcCWMPb.UpdateUsernameRequest) (*grpcCWMPb.UpdateUsernameResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("UpdateProfile - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	username := req.GetUserName()
	if len(username) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid username")
	}

	user, err := dao.GetUserDAO().FindByUserName(ctx, username)
	if user != nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("User existed: %v", username))
	}

	err = dao.GetUserDAO().DeleteFromCache(ctx, grpcSession.User)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	updateFields := primitive.M{}
	updateFields["username"] = username
	update := primitive.M{"$set": updateFields}
	user, err = dao.GetUserDAO().UpdateByPhoneFull(ctx, grpcSession.User.PhoneFull, update, []interface{}{}, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	go func(userOTT model.User) {
		err := sv.SendNotifyUpdateContactOTT(userOTT)
		if err != nil {
			log.Println("Cannot SendNotifyUpdateContactOTT", err)
		}
	}(*user)

	return &grpcCWMPb.UpdateUsernameResponse{
		UserName: user.Username,
	}, nil
}

func (sv *CWMGRPCService) UpdatePushToken(ctx context.Context, req *grpcCWMPb.UpdatePushTokenRequest) (*grpcCWMPb.UpdatePushTokenResponse, error) {
	grpcSession, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("UpdatePushToken - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}
	log.Println("call UpdatePushToken")

	pushTokenInfo := req.GetPushTokenInfo()

	idx := slices.IndexFunc(grpcSession.User.Sessions, func(u model.UserSession) bool { return u.SessionId == grpcSession.SessionId })
	if idx < 0 {
		return nil, status.Errorf(codes.Internal, fmt.Sprint("Cannot get UserSession"))
	}

	user := grpcSession.User
	userSession := user.Sessions[idx]

	updateFields := primitive.M{}
	if userSession.OsType == grpcCWMPb.OS_TYPE_ANDROID {
		if pushTokenInfo.GetPushTokenServiceType() != grpcCWMPb.PUSH_TOKEN_SERVICE_TYPE_FCM {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprint("Invalid push token service"))
		}
		updateFields["sessions.$[updateSession].pushtokenID"] = pushTokenInfo.GetPushtokenID()
	} else if userSession.OsType == grpcCWMPb.OS_TYPE_IOS {
		if pushTokenInfo.GetPushTokenServiceType() != grpcCWMPb.PUSH_TOKEN_SERVICE_TYPE_APNS_VOIP &&
			pushTokenInfo.GetPushTokenServiceType() != grpcCWMPb.PUSH_TOKEN_SERVICE_TYPE_APNS_REMOTE {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprint("Invalid push token service"))
		}

		if pushTokenInfo.GetPushTokenServiceType() != grpcCWMPb.PUSH_TOKEN_SERVICE_TYPE_APNS_REMOTE {
			updateFields["sessions.$[updateSession].pushtokenID"] = pushTokenInfo.GetPushtokenID()
		} else {
			updateFields["sessions.$[updateSession].secondaryPushtokenID"] = pushTokenInfo.GetPushtokenID()
		}
	}

	updateFields["sessions.$[updateSession].bundleId"] = pushTokenInfo.GetBundleid()
	updateFields["sessions.$[updateSession].appId"] = pushTokenInfo.GetAppid()

	arrayFilter := primitive.M{}
	arrayFilter["updateSession.sessionId"] = userSession.SessionId
	arrayFilters := []interface{}{arrayFilter}

	update := primitive.M{"$set": updateFields}
	user, err := dao.GetUserDAO().UpdateByPhoneFull(context.Background(), user.PhoneFull, update, arrayFilters, false)

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	return &grpcCWMPb.UpdatePushTokenResponse{
		PushTokenInfo: pushTokenInfo,
	}, nil
}

func (sv *CWMGRPCService) SearchByUsername(ctx context.Context, req *grpcCWMPb.SearchByUsernameRequest) (*grpcCWMPb.SearchByUsernameResponse, error) {
	_, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("SearchByUsername - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}
	username := req.GetUserName()
	searchUserInfos := []*grpcCWMPb.SearchUserInfo{}

	users, err := dao.GetUserDAO().SearchByUserName(ctx, username)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) { //NotFound user
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}

	if users != nil && len(users) > 0 {
		for _, user := range users {

			searchUserInfo := grpcCWMPb.SearchUserInfo{
				PhoneFull:  user.PhoneFull,
				UserId:     user.ID.Hex(),
				Username:   user.Username,
				UserAvatar: user.Avatar,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
			}
			searchUserInfos = append(searchUserInfos, &searchUserInfo)
		}
	}

	return &grpcCWMPb.SearchByUsernameResponse{
		SearchUserInfos: searchUserInfos,
	}, nil
}

func (sv *CWMGRPCService) SearchByPhoneFull(ctx context.Context, req *grpcCWMPb.SearchByPhoneFullRequest) (*grpcCWMPb.SearchByPhoneFullResponse, error) {
	_, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("SearchByPhoneFull - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}
	phoneFull := req.GetPhoneFull()
	searchUserInfos := []*grpcCWMPb.SearchUserInfo{}

	users, err := dao.GetUserDAO().SearchByPhoneFull(ctx, phoneFull)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) { //NotFound user
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
		}
	}

	if users != nil && len(users) > 0 {
		for _, user := range users {

			searchUserInfo := grpcCWMPb.SearchUserInfo{
				PhoneFull:  user.PhoneFull,
				UserId:     user.ID.Hex(),
				Username:   user.Username,
				UserAvatar: user.Avatar,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
			}
			searchUserInfos = append(searchUserInfos, &searchUserInfo)
		}
	}

	return &grpcCWMPb.SearchByPhoneFullResponse{
		SearchUserInfos: searchUserInfos,
	}, nil
}

func (sv *CWMGRPCService) FindByListPhoneFull(ctx context.Context, req *grpcCWMPb.FindByListPhoneFullRequest) (*grpcCWMPb.FindByListPhoneFullResponse, error) {
	_, ok := ctx.Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("FindByListPhoneFull - can not cast GrpcSession")
		return nil, status.Errorf(codes.PermissionDenied, "Invalid GrpcSession")
	}

	listPhoneFull := req.GetPhoneFulls()
	users, err := dao.GetUserDAO().FindByListPhoneFull(ctx, listPhoneFull)
	if err != nil {
		log.Println("Internal err:", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
	}

	var searchUserInfos []*grpcCWMPb.SearchUserInfo
	for _, participantUser := range users {
		searchUserInfo := &grpcCWMPb.SearchUserInfo{
			PhoneFull:  participantUser.PhoneFull,
			UserId:     participantUser.ID.Hex(),
			Username:   participantUser.Username,
			UserAvatar: participantUser.Avatar,
			FirstName:  participantUser.FirstName,
			LastName:   participantUser.LastName,
		}
		searchUserInfos = append(searchUserInfos, searchUserInfo)
	}

	return &grpcCWMPb.FindByListPhoneFullResponse{
		SearchUserInfos: searchUserInfos,
	}, nil
}
