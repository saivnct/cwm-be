package appgrpc

import (
	"sol.go/cwm/proto/grpcCWMPb"
)

type CWMGRPCService struct {
	grpcCWMPb.UnimplementedCWMServiceServer
}
