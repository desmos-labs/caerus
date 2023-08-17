// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: caerus/grants/v1/service.proto

package grants

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
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

type RequestFeeAllowanceRequest struct {
	// UserDesmosAddress represents the Desmos address of the user that
	// should be granted the fee allowance
	UserDesmosAddress string `protobuf:"bytes,1,opt,name=user_desmos_address,json=userDesmosAddress,proto3" json:"user_desmos_address,omitempty"`
	// Allowance represents the fee allowance that will be granted to the user.
	// IT can be any allowance type that implements AllowanceI
	Allowance *types.Any `protobuf:"bytes,2,opt,name=allowance,proto3" json:"allowance,omitempty"`
}

func (m *RequestFeeAllowanceRequest) Reset()         { *m = RequestFeeAllowanceRequest{} }
func (m *RequestFeeAllowanceRequest) String() string { return proto.CompactTextString(m) }
func (*RequestFeeAllowanceRequest) ProtoMessage()    {}
func (*RequestFeeAllowanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3408ad1c94df29ed, []int{0}
}
func (m *RequestFeeAllowanceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RequestFeeAllowanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RequestFeeAllowanceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RequestFeeAllowanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestFeeAllowanceRequest.Merge(m, src)
}
func (m *RequestFeeAllowanceRequest) XXX_Size() int {
	return m.Size()
}
func (m *RequestFeeAllowanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestFeeAllowanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RequestFeeAllowanceRequest proto.InternalMessageInfo

func (m *RequestFeeAllowanceRequest) GetUserDesmosAddress() string {
	if m != nil {
		return m.UserDesmosAddress
	}
	return ""
}

func (m *RequestFeeAllowanceRequest) GetAllowance() *types.Any {
	if m != nil {
		return m.Allowance
	}
	return nil
}

func init() {
	proto.RegisterType((*RequestFeeAllowanceRequest)(nil), "caerus.grants.v1.RequestFeeAllowanceRequest")
}

func init() { proto.RegisterFile("caerus/grants/v1/service.proto", fileDescriptor_3408ad1c94df29ed) }

var fileDescriptor_3408ad1c94df29ed = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xdf, 0x4a, 0x3a, 0x41,
	0x14, 0xc7, 0x9d, 0xdf, 0xc5, 0x0f, 0xdc, 0x08, 0x6a, 0x8d, 0xd0, 0x0d, 0x06, 0xf1, 0x22, 0x84,
	0xf2, 0x0c, 0x6b, 0x4f, 0xa0, 0x64, 0xd1, 0xad, 0xdd, 0xd9, 0x85, 0xcc, 0xae, 0xc7, 0x4d, 0x58,
	0x77, 0x6c, 0xfe, 0x6c, 0xf8, 0x16, 0x3d, 0x46, 0x0f, 0xd0, 0x43, 0x44, 0x57, 0x5e, 0x76, 0x19,
	0xfa, 0x22, 0xe1, 0xcc, 0x6e, 0x81, 0xd5, 0xe5, 0x99, 0xcf, 0x70, 0xbe, 0x7f, 0x8e, 0x47, 0x63,
	0x8e, 0xd2, 0x28, 0x96, 0x48, 0x9e, 0x69, 0xc5, 0xf2, 0x90, 0x29, 0x94, 0xf9, 0x2c, 0x46, 0x58,
	0x48, 0xa1, 0x85, 0x7f, 0xe0, 0x38, 0x38, 0x0e, 0x79, 0x18, 0x34, 0x62, 0xa1, 0xe6, 0x42, 0x8d,
	0x2d, 0x67, 0x6e, 0x70, 0x9f, 0x83, 0x46, 0x22, 0x44, 0x92, 0x22, 0xb3, 0x53, 0x64, 0xa6, 0x8c,
	0x67, 0xcb, 0x02, 0x9d, 0xec, 0x22, 0x9c, 0x2f, 0x74, 0x01, 0x5b, 0xcf, 0xc4, 0x0b, 0x86, 0xf8,
	0x60, 0x50, 0xe9, 0x2b, 0xc4, 0x5e, 0x9a, 0x8a, 0x47, 0x9e, 0xc5, 0x58, 0x3c, 0xf9, 0xe0, 0xd5,
	0x8c, 0x42, 0x39, 0x9e, 0xa0, 0x15, 0xe6, 0x93, 0x89, 0x44, 0xa5, 0xea, 0xa4, 0x49, 0xda, 0xd5,
	0xe1, 0xe1, 0x16, 0x5d, 0x5a, 0xd2, 0x73, 0xc0, 0x1f, 0x79, 0x55, 0x5e, 0xee, 0xa8, 0xff, 0x6b,
	0x92, 0xf6, 0x5e, 0xf7, 0x08, 0x9c, 0x3e, 0x94, 0xfa, 0xd0, 0xcb, 0x96, 0xfd, 0xd3, 0xb7, 0x97,
	0x4e, 0xab, 0x48, 0x30, 0x45, 0xb4, 0x19, 0x21, 0x0f, 0x23, 0xd4, 0x3c, 0x84, 0x2f, 0x1b, 0x37,
	0xc3, 0xef, 0x75, 0xdd, 0xd4, 0xdb, 0xbf, 0xb6, 0x55, 0xdc, 0xba, 0x9a, 0xfc, 0x3b, 0xaf, 0xf6,
	0x8b, 0x75, 0xff, 0x1c, 0x76, 0x8b, 0x83, 0xbf, 0x13, 0x06, 0xc7, 0x3f, 0xec, 0x0d, 0xb6, 0xf5,
	0xf4, 0x07, 0xaf, 0x6b, 0x4a, 0x56, 0x6b, 0x4a, 0x3e, 0xd6, 0x94, 0x3c, 0x6d, 0x68, 0x65, 0xb5,
	0xa1, 0x95, 0xf7, 0x0d, 0xad, 0x8c, 0xce, 0x92, 0x99, 0xbe, 0x37, 0x11, 0xc4, 0x62, 0xce, 0x5c,
	0x2f, 0x9d, 0x94, 0x47, 0x8a, 0x15, 0xe7, 0x94, 0xc2, 0x68, 0x2c, 0xaf, 0x1a, 0xfd, 0xb7, 0x6b,
	0x2f, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xef, 0x38, 0xa3, 0x7d, 0xed, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GrantsServiceClient is the client API for GrantsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GrantsServiceClient interface {
	RequestFeeAllowance(ctx context.Context, in *RequestFeeAllowanceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type grantsServiceClient struct {
	cc grpc1.ClientConn
}

func NewGrantsServiceClient(cc grpc1.ClientConn) GrantsServiceClient {
	return &grantsServiceClient{cc}
}

func (c *grantsServiceClient) RequestFeeAllowance(ctx context.Context, in *RequestFeeAllowanceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/caerus.grants.v1.GrantsService/RequestFeeAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrantsServiceServer is the server API for GrantsService service.
type GrantsServiceServer interface {
	RequestFeeAllowance(context.Context, *RequestFeeAllowanceRequest) (*emptypb.Empty, error)
}

// UnimplementedGrantsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGrantsServiceServer struct {
}

func (*UnimplementedGrantsServiceServer) RequestFeeAllowance(ctx context.Context, req *RequestFeeAllowanceRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestFeeAllowance not implemented")
}

func RegisterGrantsServiceServer(s grpc1.Server, srv GrantsServiceServer) {
	s.RegisterService(&_GrantsService_serviceDesc, srv)
}

func _GrantsService_RequestFeeAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestFeeAllowanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrantsServiceServer).RequestFeeAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/caerus.grants.v1.GrantsService/RequestFeeAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrantsServiceServer).RequestFeeAllowance(ctx, req.(*RequestFeeAllowanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GrantsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "caerus.grants.v1.GrantsService",
	HandlerType: (*GrantsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestFeeAllowance",
			Handler:    _GrantsService_RequestFeeAllowance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "caerus/grants/v1/service.proto",
}

func (m *RequestFeeAllowanceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RequestFeeAllowanceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RequestFeeAllowanceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Allowance != nil {
		{
			size, err := m.Allowance.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintService(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.UserDesmosAddress) > 0 {
		i -= len(m.UserDesmosAddress)
		copy(dAtA[i:], m.UserDesmosAddress)
		i = encodeVarintService(dAtA, i, uint64(len(m.UserDesmosAddress)))
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
func (m *RequestFeeAllowanceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UserDesmosAddress)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	if m.Allowance != nil {
		l = m.Allowance.Size()
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
func (m *RequestFeeAllowanceRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: RequestFeeAllowanceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RequestFeeAllowanceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserDesmosAddress", wireType)
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
			m.UserDesmosAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allowance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Allowance == nil {
				m.Allowance = &types.Any{}
			}
			if err := m.Allowance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
