// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: caerus/applications/v1/service.proto

package applications

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type RegisterNotificationTokenRequest struct {
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (m *RegisterNotificationTokenRequest) Reset()         { *m = RegisterNotificationTokenRequest{} }
func (m *RegisterNotificationTokenRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterNotificationTokenRequest) ProtoMessage()    {}
func (*RegisterNotificationTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0f1298cb9f15fac, []int{0}
}
func (m *RegisterNotificationTokenRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegisterNotificationTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegisterNotificationTokenRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegisterNotificationTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterNotificationTokenRequest.Merge(m, src)
}
func (m *RegisterNotificationTokenRequest) XXX_Size() int {
	return m.Size()
}
func (m *RegisterNotificationTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterNotificationTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterNotificationTokenRequest proto.InternalMessageInfo

func (m *RegisterNotificationTokenRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type DeleteAppRequest struct {
	AppId string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
}

func (m *DeleteAppRequest) Reset()         { *m = DeleteAppRequest{} }
func (m *DeleteAppRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteAppRequest) ProtoMessage()    {}
func (*DeleteAppRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0f1298cb9f15fac, []int{1}
}
func (m *DeleteAppRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeleteAppRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeleteAppRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeleteAppRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAppRequest.Merge(m, src)
}
func (m *DeleteAppRequest) XXX_Size() int {
	return m.Size()
}
func (m *DeleteAppRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAppRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAppRequest proto.InternalMessageInfo

func (m *DeleteAppRequest) GetAppId() string {
	if m != nil {
		return m.AppId
	}
	return ""
}

func init() {
	proto.RegisterType((*RegisterNotificationTokenRequest)(nil), "caerus.applications.v1.RegisterNotificationTokenRequest")
	proto.RegisterType((*DeleteAppRequest)(nil), "caerus.applications.v1.DeleteAppRequest")
}

func init() {
	proto.RegisterFile("caerus/applications/v1/service.proto", fileDescriptor_f0f1298cb9f15fac)
}

var fileDescriptor_f0f1298cb9f15fac = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xbf, 0x4b, 0xc3, 0x40,
	0x14, 0xc7, 0x7b, 0x83, 0x85, 0xde, 0x24, 0x87, 0x16, 0xad, 0x70, 0x94, 0xe2, 0x50, 0x07, 0xef,
	0x88, 0x2e, 0x5d, 0x2b, 0x3a, 0x88, 0xe8, 0x10, 0x9d, 0x5c, 0xe4, 0x92, 0xbc, 0xc6, 0xc3, 0xa4,
	0x77, 0xe6, 0x5e, 0x02, 0xfe, 0x17, 0xfe, 0x59, 0x8e, 0xc5, 0xc9, 0x51, 0x92, 0x7f, 0x44, 0xda,
	0x34, 0xfe, 0xc2, 0xe0, 0xf8, 0xe0, 0xf3, 0xbe, 0xef, 0xfb, 0xe1, 0xd1, 0xfd, 0x50, 0x41, 0x96,
	0x3b, 0xa9, 0xac, 0x4d, 0x74, 0xa8, 0x50, 0x9b, 0xb9, 0x93, 0x85, 0x27, 0x1d, 0x64, 0x85, 0x0e,
	0x41, 0xd8, 0xcc, 0xa0, 0x61, 0xfd, 0x9a, 0x12, 0xdf, 0x29, 0x51, 0x78, 0x83, 0xbd, 0xd8, 0x98,
	0x38, 0x01, 0xb9, 0xa2, 0x82, 0x7c, 0x26, 0x21, 0xb5, 0xf8, 0x54, 0x2f, 0x8d, 0x26, 0x74, 0xe8,
	0x43, 0xac, 0x1d, 0x42, 0x76, 0x65, 0x50, 0xcf, 0xd6, 0x8b, 0x37, 0xe6, 0x01, 0xe6, 0x3e, 0x3c,
	0xe6, 0xe0, 0x90, 0x6d, 0xd1, 0x0d, 0x5c, 0xce, 0x3b, 0x64, 0x48, 0xc6, 0x3d, 0xbf, 0x1e, 0x46,
	0x07, 0x74, 0xf3, 0x14, 0x12, 0x40, 0x98, 0x5a, 0xdb, 0x90, 0xdb, 0xb4, 0xab, 0xac, 0xbd, 0xd3,
	0x51, 0x83, 0x2a, 0x6b, 0xcf, 0xa3, 0xa3, 0x57, 0x42, 0xd9, 0xf4, 0xab, 0xd5, 0x75, 0x5d, 0x9b,
	0xa5, 0x74, 0xb7, 0xf5, 0x36, 0x9b, 0x88, 0xbf, 0x75, 0xc4, 0x7f, 0x75, 0x07, 0x7d, 0x51, 0x0b,
	0x8b, 0x46, 0x58, 0x9c, 0x2d, 0x85, 0xd9, 0x25, 0xed, 0x7d, 0x16, 0x66, 0xe3, 0xb6, 0xf8, 0xdf,
	0x4e, 0x6d, 0x71, 0x27, 0x17, 0x2f, 0x25, 0x27, 0x8b, 0x92, 0x93, 0xf7, 0x92, 0x93, 0xe7, 0x8a,
	0x77, 0x16, 0x15, 0xef, 0xbc, 0x55, 0xbc, 0x73, 0xeb, 0xc5, 0x1a, 0xef, 0xf3, 0x40, 0x84, 0x26,
	0x95, 0x11, 0xb8, 0xd4, 0xb8, 0xc3, 0x44, 0x05, 0x4e, 0xae, 0xbf, 0x98, 0x99, 0x1c, 0xe1, 0xe7,
	0x33, 0x83, 0xee, 0x2a, 0xfc, 0xf8, 0x23, 0x00, 0x00, 0xff, 0xff, 0x7d, 0xfb, 0x8f, 0xb5, 0xea,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApplicationServiceClient is the client API for ApplicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApplicationServiceClient interface {
	// RegisterNotificationToken allows to register a notification token for a
	// given application
	RegisterNotificationToken(ctx context.Context, in *RegisterNotificationTokenRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteApp allows to delete an application
	DeleteApp(ctx context.Context, in *DeleteAppRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type applicationServiceClient struct {
	cc grpc1.ClientConn
}

func NewApplicationServiceClient(cc grpc1.ClientConn) ApplicationServiceClient {
	return &applicationServiceClient{cc}
}

func (c *applicationServiceClient) RegisterNotificationToken(ctx context.Context, in *RegisterNotificationTokenRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/caerus.applications.v1.ApplicationService/RegisterNotificationToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationServiceClient) DeleteApp(ctx context.Context, in *DeleteAppRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/caerus.applications.v1.ApplicationService/DeleteApp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApplicationServiceServer is the server API for ApplicationService service.
type ApplicationServiceServer interface {
	// RegisterNotificationToken allows to register a notification token for a
	// given application
	RegisterNotificationToken(context.Context, *RegisterNotificationTokenRequest) (*emptypb.Empty, error)
	// DeleteApp allows to delete an application
	DeleteApp(context.Context, *DeleteAppRequest) (*emptypb.Empty, error)
}

// UnimplementedApplicationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedApplicationServiceServer struct {
}

func (*UnimplementedApplicationServiceServer) RegisterNotificationToken(ctx context.Context, req *RegisterNotificationTokenRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNotificationToken not implemented")
}
func (*UnimplementedApplicationServiceServer) DeleteApp(ctx context.Context, req *DeleteAppRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteApp not implemented")
}

func RegisterApplicationServiceServer(s grpc1.Server, srv ApplicationServiceServer) {
	s.RegisterService(&_ApplicationService_serviceDesc, srv)
}

func _ApplicationService_RegisterNotificationToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterNotificationTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).RegisterNotificationToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/caerus.applications.v1.ApplicationService/RegisterNotificationToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).RegisterNotificationToken(ctx, req.(*RegisterNotificationTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApplicationService_DeleteApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServiceServer).DeleteApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/caerus.applications.v1.ApplicationService/DeleteApp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServiceServer).DeleteApp(ctx, req.(*DeleteAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApplicationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "caerus.applications.v1.ApplicationService",
	HandlerType: (*ApplicationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNotificationToken",
			Handler:    _ApplicationService_RegisterNotificationToken_Handler,
		},
		{
			MethodName: "DeleteApp",
			Handler:    _ApplicationService_DeleteApp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "caerus/applications/v1/service.proto",
}

func (m *RegisterNotificationTokenRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegisterNotificationTokenRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegisterNotificationTokenRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintService(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DeleteAppRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeleteAppRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeleteAppRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AppId) > 0 {
		i -= len(m.AppId)
		copy(dAtA[i:], m.AppId)
		i = encodeVarintService(dAtA, i, uint64(len(m.AppId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintService(dAtA []byte, offset int, v uint64) int {
	offset -= sovService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RegisterNotificationTokenRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	return n
}

func (m *DeleteAppRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AppId)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	return n
}

func sovService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozService(x uint64) (n int) {
	return sovService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RegisterNotificationTokenRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RegisterNotificationTokenRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegisterNotificationTokenRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DeleteAppRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DeleteAppRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeleteAppRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AppId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupService = fmt.Errorf("proto: unexpected end of group")
)
