// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: grpc/cwm-sv.proto

package grpcCWMPb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_grpc_cwm_sv_proto protoreflect.FileDescriptor

var file_grpc_cwm_sv_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x77, 0x6d, 0x2d, 0x73, 0x76, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x1a, 0x1d,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x77, 0x6d, 0x2d, 0x72, 0x71, 0x2d, 0x72, 0x65, 0x73, 0x2d,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x63, 0x77, 0x6d, 0x2d, 0x72, 0x71, 0x2d, 0x72, 0x65, 0x73, 0x2d, 0x74,
	0x68, 0x72, 0x65, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x63, 0x77, 0x6d, 0x2d, 0x72, 0x71, 0x2d, 0x72, 0x65, 0x73, 0x2d, 0x6d, 0x73, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x92, 0x16, 0x0a, 0x0a, 0x43, 0x57, 0x4d, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x1e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x10, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57,
	0x4d, 0x50, 0x62, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3a, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x17, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x52, 0x0a, 0x0b,
	0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x12, 0x1d, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x74, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01,
	0x12, 0x52, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x1f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d,
	0x50, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43,
	0x57, 0x4d, 0x50, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x10, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a, 0x11, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x42, 0x79, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x12, 0x23, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x42, 0x79, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x42, 0x79, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x64, 0x0a, 0x13, 0x46, 0x69, 0x6e, 0x64,
	0x42, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x12,
	0x25, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x42, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d,
	0x50, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x68, 0x6f,
	0x6e, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58,
	0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x12, 0x23, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x67, 0x0a, 0x14, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x26, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43,
	0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54,
	0x68, 0x72, 0x65, 0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x6a, 0x0a, 0x15, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61,
	0x64, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x76, 0x0a,
	0x19, 0x41, 0x64, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x50,
	0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x12, 0x2b, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54,
	0x68, 0x72, 0x65, 0x61, 0x64, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57,
	0x4d, 0x50, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65,
	0x61, 0x64, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7f, 0x0a, 0x1c, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x6e, 0x74, 0x12, 0x2e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50,
	0x62, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50,
	0x62, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x70, 0x0a, 0x17, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x12, 0x29, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x50, 0x72,
	0x6f, 0x6d, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6d, 0x0a, 0x16, 0x52, 0x65, 0x76, 0x6f,
	0x6b, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x12, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x52,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x10, 0x4c, 0x65, 0x61, 0x76, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x12, 0x22, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x23, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x4c, 0x65, 0x61, 0x76,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x76, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6e,
	0x64, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68, 0x72, 0x65, 0x61,
	0x64, 0x12, 0x2b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x41, 0x6e, 0x64, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x6e, 0x64, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x68,
	0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x0e,
	0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x12, 0x20,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x69,
	0x61, 0x6c, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x49, 0x6e, 0x69,
	0x74, 0x69, 0x61, 0x6c, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x6e, 0x0a, 0x15, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x41, 0x6c, 0x6c, 0x55, 0x6e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d, 0x73, 0x67,
	0x12, 0x27, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x41, 0x6c, 0x6c, 0x55, 0x6e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d,
	0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x41, 0x6c, 0x6c, 0x55, 0x6e,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x68, 0x0a, 0x13, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x4f, 0x6c, 0x64, 0x4d, 0x73, 0x67, 0x4f, 0x66, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x12, 0x25,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x4f, 0x6c, 0x64, 0x4d, 0x73, 0x67, 0x4f, 0x66, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50,
	0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x4f, 0x6c, 0x64, 0x4d, 0x73, 0x67, 0x4f, 0x66, 0x54,
	0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30,
	0x01, 0x12, 0x40, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x19, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57,
	0x4d, 0x50, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x64, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x73, 0x12, 0x25, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x26, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x4d, 0x73, 0x67,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x12, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x4f, 0x66, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x12,
	0x24, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x4f, 0x66, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50,
	0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x73, 0x4f, 0x66, 0x54, 0x68,
	0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x64, 0x0a, 0x13,
	0x43, 0x6c, 0x65, 0x61, 0x72, 0x41, 0x6c, 0x6c, 0x4d, 0x73, 0x67, 0x4f, 0x66, 0x54, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x12, 0x25, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e,
	0x43, 0x6c, 0x65, 0x61, 0x72, 0x41, 0x6c, 0x6c, 0x4d, 0x73, 0x67, 0x4f, 0x66, 0x54, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x41, 0x6c, 0x6c, 0x4d,
	0x73, 0x67, 0x4f, 0x66, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x5b, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6f, 0x6c, 0x6f,
	0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x12, 0x22, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d,
	0x50, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6f, 0x6c, 0x6f, 0x54, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6f, 0x6c,
	0x6f, 0x54, 0x68, 0x72, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x59, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73,
	0x67, 0x12, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x5f, 0x0a, 0x10, 0x44, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73, 0x67, 0x12, 0x22,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c,
	0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x23, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x2e, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x1c, 0x5a, 0x1a, 0x73,
	0x6f, 0x6c, 0x2e, 0x67, 0x6f, 0x2f, 0x63, 0x77, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x43, 0x57, 0x4d, 0x50, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_grpc_cwm_sv_proto_goTypes = []interface{}{
	(*CreatAccountRequest)(nil),                  // 0: grpcCWMPb.CreatAccountRequest
	(*VerifyAuthencodeRequest)(nil),              // 1: grpcCWMPb.VerifyAuthencodeRequest
	(*LoginRequest)(nil),                         // 2: grpcCWMPb.LoginRequest
	(*SyncContactRequest)(nil),                   // 3: grpcCWMPb.SyncContactRequest
	(*UpdateProfileRequest)(nil),                 // 4: grpcCWMPb.UpdateProfileRequest
	(*UpdateUsernameRequest)(nil),                // 5: grpcCWMPb.UpdateUsernameRequest
	(*SearchByUsernameRequest)(nil),              // 6: grpcCWMPb.SearchByUsernameRequest
	(*SearchByPhoneFullRequest)(nil),             // 7: grpcCWMPb.SearchByPhoneFullRequest
	(*FindByListPhoneFullRequest)(nil),           // 8: grpcCWMPb.FindByListPhoneFullRequest
	(*UpdatePushTokenRequest)(nil),               // 9: grpcCWMPb.UpdatePushTokenRequest
	(*CreateGroupThreadRequest)(nil),             // 10: grpcCWMPb.CreateGroupThreadRequest
	(*CheckGroupThreadInfoRequest)(nil),          // 11: grpcCWMPb.CheckGroupThreadInfoRequest
	(*ChangeGroupThreadNameRequest)(nil),         // 12: grpcCWMPb.ChangeGroupThreadNameRequest
	(*AddGroupThreadParticipantRequest)(nil),     // 13: grpcCWMPb.AddGroupThreadParticipantRequest
	(*RemoveGroupThreadParticipantRequest)(nil),  // 14: grpcCWMPb.RemoveGroupThreadParticipantRequest
	(*PromoteGroupThreadAdminRequest)(nil),       // 15: grpcCWMPb.PromoteGroupThreadAdminRequest
	(*RevokeGroupThreadAdminRequest)(nil),        // 16: grpcCWMPb.RevokeGroupThreadAdminRequest
	(*LeaveGroupThreadRequest)(nil),              // 17: grpcCWMPb.LeaveGroupThreadRequest
	(*DeleteAndLeaveGroupThreadRequest)(nil),     // 18: grpcCWMPb.DeleteAndLeaveGroupThreadRequest
	(*InitialSyncMsgRequest)(nil),                // 19: grpcCWMPb.InitialSyncMsgRequest
	(*FetchAllUnreceivedMsgRequest)(nil),         // 20: grpcCWMPb.FetchAllUnreceivedMsgRequest
	(*FetchOldMsgOfThreadRequest)(nil),           // 21: grpcCWMPb.FetchOldMsgOfThreadRequest
	(*SendMsgRequest)(nil),                       // 22: grpcCWMPb.SendMsgRequest
	(*ConfirmReceivedMsgsRequest)(nil),           // 23: grpcCWMPb.ConfirmReceivedMsgsRequest
	(*DeleteMsgsOfThreadRequest)(nil),            // 24: grpcCWMPb.DeleteMsgsOfThreadRequest
	(*ClearAllMsgOfThreadRequest)(nil),           // 25: grpcCWMPb.ClearAllMsgOfThreadRequest
	(*DeleteSoloThreadRequest)(nil),              // 26: grpcCWMPb.DeleteSoloThreadRequest
	(*UploadMediaMsgRequest)(nil),                // 27: grpcCWMPb.UploadMediaMsgRequest
	(*DownloadMediaMsgRequest)(nil),              // 28: grpcCWMPb.DownloadMediaMsgRequest
	(*CreatAccountResponse)(nil),                 // 29: grpcCWMPb.CreatAccountResponse
	(*VerifyAuthencodeResponse)(nil),             // 30: grpcCWMPb.VerifyAuthencodeResponse
	(*LoginResponse)(nil),                        // 31: grpcCWMPb.LoginResponse
	(*SyncContactResponse)(nil),                  // 32: grpcCWMPb.SyncContactResponse
	(*UpdateProfileResponse)(nil),                // 33: grpcCWMPb.UpdateProfileResponse
	(*UpdateUsernameResponse)(nil),               // 34: grpcCWMPb.UpdateUsernameResponse
	(*SearchByUsernameResponse)(nil),             // 35: grpcCWMPb.SearchByUsernameResponse
	(*SearchByPhoneFullResponse)(nil),            // 36: grpcCWMPb.SearchByPhoneFullResponse
	(*FindByListPhoneFullResponse)(nil),          // 37: grpcCWMPb.FindByListPhoneFullResponse
	(*UpdatePushTokenResponse)(nil),              // 38: grpcCWMPb.UpdatePushTokenResponse
	(*CreateGroupThreadResponse)(nil),            // 39: grpcCWMPb.CreateGroupThreadResponse
	(*CheckGroupThreadInfoResponse)(nil),         // 40: grpcCWMPb.CheckGroupThreadInfoResponse
	(*ChangeGroupThreadNameResponse)(nil),        // 41: grpcCWMPb.ChangeGroupThreadNameResponse
	(*AddGroupThreadParticipantResponse)(nil),    // 42: grpcCWMPb.AddGroupThreadParticipantResponse
	(*RemoveGroupThreadParticipantResponse)(nil), // 43: grpcCWMPb.RemoveGroupThreadParticipantResponse
	(*PromoteGroupThreadAdminResponse)(nil),      // 44: grpcCWMPb.PromoteGroupThreadAdminResponse
	(*RevokeGroupThreadAdminResponse)(nil),       // 45: grpcCWMPb.RevokeGroupThreadAdminResponse
	(*LeaveGroupThreadResponse)(nil),             // 46: grpcCWMPb.LeaveGroupThreadResponse
	(*DeleteAndLeaveGroupThreadResponse)(nil),    // 47: grpcCWMPb.DeleteAndLeaveGroupThreadResponse
	(*InitialSyncMsgResponse)(nil),               // 48: grpcCWMPb.InitialSyncMsgResponse
	(*FetchAllUnreceivedMsgResponse)(nil),        // 49: grpcCWMPb.FetchAllUnreceivedMsgResponse
	(*FetchOldMsgOfThreadResponse)(nil),          // 50: grpcCWMPb.FetchOldMsgOfThreadResponse
	(*SendMsgResponse)(nil),                      // 51: grpcCWMPb.SendMsgResponse
	(*ConfirmReceivedMsgsResponse)(nil),          // 52: grpcCWMPb.ConfirmReceivedMsgsResponse
	(*DeleteMsgsOfThreadResponse)(nil),           // 53: grpcCWMPb.DeleteMsgsOfThreadResponse
	(*ClearAllMsgOfThreadResponse)(nil),          // 54: grpcCWMPb.ClearAllMsgOfThreadResponse
	(*DeleteSoloThreadResponse)(nil),             // 55: grpcCWMPb.DeleteSoloThreadResponse
	(*UploadMediaMsgResponse)(nil),               // 56: grpcCWMPb.UploadMediaMsgResponse
	(*DownloadMediaMsgResponse)(nil),             // 57: grpcCWMPb.DownloadMediaMsgResponse
}
var file_grpc_cwm_sv_proto_depIdxs = []int32{
	0,  // 0: grpcCWMPb.CWMService.CreatUser:input_type -> grpcCWMPb.CreatAccountRequest
	1,  // 1: grpcCWMPb.CWMService.VerifyAuthencode:input_type -> grpcCWMPb.VerifyAuthencodeRequest
	2,  // 2: grpcCWMPb.CWMService.Login:input_type -> grpcCWMPb.LoginRequest
	3,  // 3: grpcCWMPb.CWMService.SyncContact:input_type -> grpcCWMPb.SyncContactRequest
	4,  // 4: grpcCWMPb.CWMService.UpdateProfile:input_type -> grpcCWMPb.UpdateProfileRequest
	5,  // 5: grpcCWMPb.CWMService.UpdateUsername:input_type -> grpcCWMPb.UpdateUsernameRequest
	6,  // 6: grpcCWMPb.CWMService.SearchByUsername:input_type -> grpcCWMPb.SearchByUsernameRequest
	7,  // 7: grpcCWMPb.CWMService.SearchByPhoneFull:input_type -> grpcCWMPb.SearchByPhoneFullRequest
	8,  // 8: grpcCWMPb.CWMService.FindByListPhoneFull:input_type -> grpcCWMPb.FindByListPhoneFullRequest
	9,  // 9: grpcCWMPb.CWMService.UpdatePushToken:input_type -> grpcCWMPb.UpdatePushTokenRequest
	10, // 10: grpcCWMPb.CWMService.CreateGroupThread:input_type -> grpcCWMPb.CreateGroupThreadRequest
	11, // 11: grpcCWMPb.CWMService.CheckGroupThreadInfo:input_type -> grpcCWMPb.CheckGroupThreadInfoRequest
	12, // 12: grpcCWMPb.CWMService.ChangeGroupThreadName:input_type -> grpcCWMPb.ChangeGroupThreadNameRequest
	13, // 13: grpcCWMPb.CWMService.AddGroupThreadParticipant:input_type -> grpcCWMPb.AddGroupThreadParticipantRequest
	14, // 14: grpcCWMPb.CWMService.RemoveGroupThreadParticipant:input_type -> grpcCWMPb.RemoveGroupThreadParticipantRequest
	15, // 15: grpcCWMPb.CWMService.PromoteGroupThreadAdmin:input_type -> grpcCWMPb.PromoteGroupThreadAdminRequest
	16, // 16: grpcCWMPb.CWMService.RevokeGroupThreadAdmin:input_type -> grpcCWMPb.RevokeGroupThreadAdminRequest
	17, // 17: grpcCWMPb.CWMService.LeaveGroupThread:input_type -> grpcCWMPb.LeaveGroupThreadRequest
	18, // 18: grpcCWMPb.CWMService.DeleteAndLeaveGroupThread:input_type -> grpcCWMPb.DeleteAndLeaveGroupThreadRequest
	19, // 19: grpcCWMPb.CWMService.InitialSyncMsg:input_type -> grpcCWMPb.InitialSyncMsgRequest
	20, // 20: grpcCWMPb.CWMService.FetchAllUnreceivedMsg:input_type -> grpcCWMPb.FetchAllUnreceivedMsgRequest
	21, // 21: grpcCWMPb.CWMService.FetchOldMsgOfThread:input_type -> grpcCWMPb.FetchOldMsgOfThreadRequest
	22, // 22: grpcCWMPb.CWMService.SendMsg:input_type -> grpcCWMPb.SendMsgRequest
	23, // 23: grpcCWMPb.CWMService.ConfirmReceivedMsgs:input_type -> grpcCWMPb.ConfirmReceivedMsgsRequest
	24, // 24: grpcCWMPb.CWMService.DeleteMsgsOfThread:input_type -> grpcCWMPb.DeleteMsgsOfThreadRequest
	25, // 25: grpcCWMPb.CWMService.ClearAllMsgOfThread:input_type -> grpcCWMPb.ClearAllMsgOfThreadRequest
	26, // 26: grpcCWMPb.CWMService.DeleteSoloThread:input_type -> grpcCWMPb.DeleteSoloThreadRequest
	27, // 27: grpcCWMPb.CWMService.UploadMediaMsg:input_type -> grpcCWMPb.UploadMediaMsgRequest
	28, // 28: grpcCWMPb.CWMService.DownloadMediaMsg:input_type -> grpcCWMPb.DownloadMediaMsgRequest
	29, // 29: grpcCWMPb.CWMService.CreatUser:output_type -> grpcCWMPb.CreatAccountResponse
	30, // 30: grpcCWMPb.CWMService.VerifyAuthencode:output_type -> grpcCWMPb.VerifyAuthencodeResponse
	31, // 31: grpcCWMPb.CWMService.Login:output_type -> grpcCWMPb.LoginResponse
	32, // 32: grpcCWMPb.CWMService.SyncContact:output_type -> grpcCWMPb.SyncContactResponse
	33, // 33: grpcCWMPb.CWMService.UpdateProfile:output_type -> grpcCWMPb.UpdateProfileResponse
	34, // 34: grpcCWMPb.CWMService.UpdateUsername:output_type -> grpcCWMPb.UpdateUsernameResponse
	35, // 35: grpcCWMPb.CWMService.SearchByUsername:output_type -> grpcCWMPb.SearchByUsernameResponse
	36, // 36: grpcCWMPb.CWMService.SearchByPhoneFull:output_type -> grpcCWMPb.SearchByPhoneFullResponse
	37, // 37: grpcCWMPb.CWMService.FindByListPhoneFull:output_type -> grpcCWMPb.FindByListPhoneFullResponse
	38, // 38: grpcCWMPb.CWMService.UpdatePushToken:output_type -> grpcCWMPb.UpdatePushTokenResponse
	39, // 39: grpcCWMPb.CWMService.CreateGroupThread:output_type -> grpcCWMPb.CreateGroupThreadResponse
	40, // 40: grpcCWMPb.CWMService.CheckGroupThreadInfo:output_type -> grpcCWMPb.CheckGroupThreadInfoResponse
	41, // 41: grpcCWMPb.CWMService.ChangeGroupThreadName:output_type -> grpcCWMPb.ChangeGroupThreadNameResponse
	42, // 42: grpcCWMPb.CWMService.AddGroupThreadParticipant:output_type -> grpcCWMPb.AddGroupThreadParticipantResponse
	43, // 43: grpcCWMPb.CWMService.RemoveGroupThreadParticipant:output_type -> grpcCWMPb.RemoveGroupThreadParticipantResponse
	44, // 44: grpcCWMPb.CWMService.PromoteGroupThreadAdmin:output_type -> grpcCWMPb.PromoteGroupThreadAdminResponse
	45, // 45: grpcCWMPb.CWMService.RevokeGroupThreadAdmin:output_type -> grpcCWMPb.RevokeGroupThreadAdminResponse
	46, // 46: grpcCWMPb.CWMService.LeaveGroupThread:output_type -> grpcCWMPb.LeaveGroupThreadResponse
	47, // 47: grpcCWMPb.CWMService.DeleteAndLeaveGroupThread:output_type -> grpcCWMPb.DeleteAndLeaveGroupThreadResponse
	48, // 48: grpcCWMPb.CWMService.InitialSyncMsg:output_type -> grpcCWMPb.InitialSyncMsgResponse
	49, // 49: grpcCWMPb.CWMService.FetchAllUnreceivedMsg:output_type -> grpcCWMPb.FetchAllUnreceivedMsgResponse
	50, // 50: grpcCWMPb.CWMService.FetchOldMsgOfThread:output_type -> grpcCWMPb.FetchOldMsgOfThreadResponse
	51, // 51: grpcCWMPb.CWMService.SendMsg:output_type -> grpcCWMPb.SendMsgResponse
	52, // 52: grpcCWMPb.CWMService.ConfirmReceivedMsgs:output_type -> grpcCWMPb.ConfirmReceivedMsgsResponse
	53, // 53: grpcCWMPb.CWMService.DeleteMsgsOfThread:output_type -> grpcCWMPb.DeleteMsgsOfThreadResponse
	54, // 54: grpcCWMPb.CWMService.ClearAllMsgOfThread:output_type -> grpcCWMPb.ClearAllMsgOfThreadResponse
	55, // 55: grpcCWMPb.CWMService.DeleteSoloThread:output_type -> grpcCWMPb.DeleteSoloThreadResponse
	56, // 56: grpcCWMPb.CWMService.UploadMediaMsg:output_type -> grpcCWMPb.UploadMediaMsgResponse
	57, // 57: grpcCWMPb.CWMService.DownloadMediaMsg:output_type -> grpcCWMPb.DownloadMediaMsgResponse
	29, // [29:58] is the sub-list for method output_type
	0,  // [0:29] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_cwm_sv_proto_init() }
func file_grpc_cwm_sv_proto_init() {
	if File_grpc_cwm_sv_proto != nil {
		return
	}
	file_grpc_cwm_rq_res_account_proto_init()
	file_grpc_cwm_rq_res_thread_proto_init()
	file_grpc_cwm_rq_res_msg_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_cwm_sv_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_cwm_sv_proto_goTypes,
		DependencyIndexes: file_grpc_cwm_sv_proto_depIdxs,
	}.Build()
	File_grpc_cwm_sv_proto = out.File
	file_grpc_cwm_sv_proto_rawDesc = nil
	file_grpc_cwm_sv_proto_goTypes = nil
	file_grpc_cwm_sv_proto_depIdxs = nil
}
