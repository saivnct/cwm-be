package appgrpc

import (
	"context"
	"sol.go/cwm/model"
)

type GrpcSession struct {
	User      *model.User
	SessionId string
}

type GrpcLoginCreds struct {
	Authorization string
	PhoneFull     string
}

func (c *GrpcLoginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		GRPC_CRED_AUTH:   c.Authorization,
		GRPC_CRED_CLIENT: c.PhoneFull,
	}, nil
}

func (c *GrpcLoginCreds) RequireTransportSecurity() bool {
	//If you don't need transport credentials, create the channel with `grpc.WithInsecure()`.
	//And make sure your `TokenAuth` returns false in `RequireTransportSecurity()`. Otherwise Dial will fail.
	//https://groups.google.com/g/grpc-io/c/sN30bEPNr6o?pli=1
	return false
}
