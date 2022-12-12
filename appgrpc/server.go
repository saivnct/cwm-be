package appgrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/utils"
	"strconv"
)

const (
	GRPC_CRED_AUTH   = "authorization"
	GRPC_CRED_CLIENT = "client"

	GRPC_CTX_KEY_SESSION = "GRPC_CTX_KEY_SESSION"
)

var (
	GRPCAccessDeniedErr        = status.Errorf(codes.PermissionDenied, "Access Denied")
	GRPCUnauthenticateUserdErr = status.Errorf(codes.PermissionDenied, "Invalid User")
	GRPCInvalidSessionErr      = status.Errorf(codes.PermissionDenied, "Invalid Session")
)

func StartServer(grpcPort int) (net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen grpc server: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcAuthentication, _ := strconv.ParseBool(os.Getenv("GRPC_AUTHENTICATION"))
	if grpcAuthentication {
		opts = append(opts, grpc.StreamInterceptor(streamInterceptor), grpc.UnaryInterceptor(unaryInterceptor))
	}

	tls, _ := strconv.ParseBool(os.Getenv("GRPC_TLS"))

	fmt.Println("GRPC TLS", tls)

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed to loading grpc certificate: %v", err)
		}
		//spew.Dump(creds)
		opts = append(opts, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(opts...)

	cwmGRPCService := CWMGRPCService{}
	grpcCWMPb.RegisterCWMServiceServer(grpcServer, &cwmGRPCService)

	//Register reflection service on gRPC server, so that evan CLI can interact with server without .proto files
	reflection.Register(grpcServer)

	go func() {
		log.Printf("Starting grpc server on: %v\n", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve grpc: %v", err)
		}
	}()

	return listener, grpcServer, nil
}

func authorize(ctx context.Context) (*model.User, string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//spew.Dump(md)
		if len(md[GRPC_CRED_AUTH]) > 0 && len(md[GRPC_CRED_CLIENT]) > 0 {
			jwtToken := md[GRPC_CRED_AUTH][0]
			phoneFull := md[GRPC_CRED_CLIENT][0]
			//log.Printf("authorize - %v - %v", phoneFull, jwtToken)

			if len(jwtToken) > 0 && len(phoneFull) > 0 {
				user, err := dao.GetUserDAO().FindByPhoneFull(ctx, phoneFull)
				if err != nil {
					return nil, "", GRPCUnauthenticateUserdErr
				}

				claims, err := utils.ParseJWTToken(jwtToken)

				if err != nil {
					log.Println("GRPC - Invalid jwtToken", err)

					if errors.Is(err, jwt.ErrTokenExpired) {
						log.Println("GRPC - jwtToken Expired!!!!!")
						if user.IsExpireNonce() {
							if user.ShouldRenewNonce() {
								log.Printf("GRPC - %v shouldRenewNonce!!!!!", phoneFull)
								user, err = dao.GetUserDAO().UpdateNonce(ctx, user.PhoneFull)
								if err != nil {
									return nil, "", status.Errorf(codes.Internal, fmt.Sprintf("Internal err: %v", err))
								}
							}
						}
					}

					return nil, "", status.Errorf(codes.Unauthenticated, user.Nonce)
				}

				jwtPhoneFull := claims.Subject
				jwtSessionId := claims.ID

				if user.PhoneFull != jwtPhoneFull {
					return nil, "", GRPCUnauthenticateUserdErr
				}

				idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == jwtSessionId })
				if idx < 0 {
					return nil, "", GRPCInvalidSessionErr
				}

				return user, jwtSessionId, nil
			}
		}
	}

	return nil, "", GRPCAccessDeniedErr
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch req.(type) {
	case
		*grpcCWMPb.CreatAccountRequest,
		*grpcCWMPb.VerifyAuthencodeRequest,
		*grpcCWMPb.LoginRequest:
		//log.Println("unaryInterceptor: bypass authorization check")
		return handler(ctx, req)
	default:
		//log.Println("unaryInterceptor: should authorization check")
	}

	user, sessionId, err := authorize(ctx)
	if err != nil {
		grpcDummyTestMode, _ := strconv.ParseBool(os.Getenv("GRPC_DUMMY_TEST_MODE"))
		grpcDummyUserPhoneFull := os.Getenv("GRPC_DUMMY_USER")
		grpcDummySessionId := os.Getenv("GRPC_DUMMY_SESSION")
		if grpcDummyTestMode && len(grpcDummyUserPhoneFull) > 0 && len(grpcDummySessionId) > 0 {
			log.Println("Using GRPC DummyUser", grpcDummyUserPhoneFull, grpcDummySessionId)

			user, err = dao.GetUserDAO().FindByPhoneFull(ctx, grpcDummyUserPhoneFull)
			if err != nil {
				log.Println("Using GRPC DummyUser Err", err)
				return nil, err
			}

			idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == grpcDummySessionId })
			if idx < 0 {
				log.Println("Using GRPC DummyUser Err Invalid SessionId", grpcDummyUserPhoneFull, grpcDummySessionId)
				return nil, GRPCUnauthenticateUserdErr
			}
			sessionId = grpcDummySessionId
		} else {
			//log.Println("GRPC_DUMMY_TEST_MODE is Disabled")
			return nil, err
		}
	}

	grpcSession := &GrpcSession{
		User:      user,
		SessionId: sessionId,
	}

	ctxInterceptor := context.WithValue(ctx, GRPC_CTX_KEY_SESSION, grpcSession)

	return handler(ctxInterceptor, req)
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	user, sessionId, err := authorize(stream.Context())
	if err != nil {
		grpcDummyTestMode, _ := strconv.ParseBool(os.Getenv("GRPC_DUMMY_TEST_MODE"))
		grpcDummyUserPhoneFull := os.Getenv("GRPC_DUMMY_USER")
		grpcDummySessionId := os.Getenv("GRPC_DUMMY_SESSION")
		if grpcDummyTestMode && len(grpcDummyUserPhoneFull) > 0 && len(grpcDummySessionId) > 0 {
			log.Println("Using GRPC DummyUser", grpcDummyUserPhoneFull, grpcDummySessionId)

			user, err = dao.GetUserDAO().FindByPhoneFull(stream.Context(), grpcDummyUserPhoneFull)
			if err != nil {
				log.Println("Using GRPC DummyUser Err", err)
				return err
			}

			idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == grpcDummySessionId })
			if idx < 0 {
				log.Println("Using GRPC DummyUser Err Invalid SessionId", grpcDummyUserPhoneFull, grpcDummySessionId)
				return GRPCUnauthenticateUserdErr
			}
			sessionId = grpcDummySessionId
		} else {
			//log.Println("GRPC_DUMMY_TEST_MODE is Disabled")
			return err
		}
	}

	return handler(srv, newWrappedStream(stream, user, sessionId))
}

// wrappedStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	//log.Println("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	//log.Println("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedStream) Context() context.Context {
	return w.WrappedContext
}

func newWrappedStream(stream grpc.ServerStream, user *model.User, sessionId string) grpc.ServerStream {
	grpcSession := &GrpcSession{
		User:      user,
		SessionId: sessionId,
	}
	wrappedContext := context.WithValue(stream.Context(), GRPC_CTX_KEY_SESSION, grpcSession)
	return &wrappedStream{stream, wrappedContext}
}
