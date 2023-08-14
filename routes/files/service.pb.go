// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: caerus/files/v1/service.proto

package files

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type FileChunk struct {
	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *FileChunk) Reset()         { *m = FileChunk{} }
func (m *FileChunk) String() string { return proto.CompactTextString(m) }
func (*FileChunk) ProtoMessage()    {}
func (*FileChunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef20e225e3647d4e, []int{0}
}
func (m *FileChunk) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FileChunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FileChunk.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FileChunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileChunk.Merge(m, src)
}
func (m *FileChunk) XXX_Size() int {
	return m.Size()
}
func (m *FileChunk) XXX_DiscardUnknown() {
	xxx_messageInfo_FileChunk.DiscardUnknown(m)
}

var xxx_messageInfo_FileChunk proto.InternalMessageInfo

func (m *FileChunk) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *FileChunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type UploadFileResponse struct {
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (m *UploadFileResponse) Reset()         { *m = UploadFileResponse{} }
func (m *UploadFileResponse) String() string { return proto.CompactTextString(m) }
func (*UploadFileResponse) ProtoMessage()    {}
func (*UploadFileResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef20e225e3647d4e, []int{1}
}
func (m *UploadFileResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UploadFileResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UploadFileResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UploadFileResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadFileResponse.Merge(m, src)
}
func (m *UploadFileResponse) XXX_Size() int {
	return m.Size()
}
func (m *UploadFileResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadFileResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UploadFileResponse proto.InternalMessageInfo

func (m *UploadFileResponse) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type GetFileRequest struct {
	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (m *GetFileRequest) Reset()         { *m = GetFileRequest{} }
func (m *GetFileRequest) String() string { return proto.CompactTextString(m) }
func (*GetFileRequest) ProtoMessage()    {}
func (*GetFileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef20e225e3647d4e, []int{2}
}
func (m *GetFileRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetFileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetFileRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetFileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFileRequest.Merge(m, src)
}
func (m *GetFileRequest) XXX_Size() int {
	return m.Size()
}
func (m *GetFileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFileRequest proto.InternalMessageInfo

func (m *GetFileRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func init() {
	proto.RegisterType((*FileChunk)(nil), "caerus.files.v1.FileChunk")
	proto.RegisterType((*UploadFileResponse)(nil), "caerus.files.v1.UploadFileResponse")
	proto.RegisterType((*GetFileRequest)(nil), "caerus.files.v1.GetFileRequest")
}

func init() { proto.RegisterFile("caerus/files/v1/service.proto", fileDescriptor_ef20e225e3647d4e) }

var fileDescriptor_ef20e225e3647d4e = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0x4e, 0x4c, 0x2d,
	0x2a, 0x2d, 0xd6, 0x4f, 0xcb, 0xcc, 0x49, 0x2d, 0xd6, 0x2f, 0x33, 0xd4, 0x2f, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x87, 0x48, 0xeb, 0x81, 0xa5,
	0xf5, 0xca, 0x0c, 0x95, 0x6c, 0xb8, 0x38, 0xdd, 0x32, 0x73, 0x52, 0x9d, 0x33, 0x4a, 0xf3, 0xb2,
	0x85, 0xa4, 0xb9, 0x38, 0x41, 0x12, 0xf1, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0x9c, 0x41, 0x1c, 0x20, 0x01, 0xbf, 0xc4, 0xdc, 0x54, 0x21, 0x21, 0x2e, 0x96, 0x94, 0xc4, 0x92,
	0x44, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x30, 0x5b, 0x49, 0x8d, 0x4b, 0x28, 0xb4, 0x20,
	0x27, 0x3f, 0x31, 0x05, 0x64, 0x46, 0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x90, 0x00,
	0x17, 0x73, 0x69, 0x51, 0x0e, 0xd4, 0x00, 0x10, 0x53, 0x49, 0x97, 0x8b, 0xcf, 0x3d, 0xb5, 0x04,
	0xa2, 0xa8, 0xb0, 0x34, 0xb5, 0xb8, 0x04, 0xaf, 0x55, 0x46, 0x2b, 0x19, 0xb9, 0x78, 0x40, 0x8a,
	0x8b, 0x83, 0x21, 0x8e, 0x17, 0xf2, 0xe7, 0xe2, 0x42, 0xd8, 0x23, 0x24, 0xa5, 0x87, 0xe6, 0x0b,
	0x3d, 0xb8, 0x17, 0xa4, 0x94, 0x31, 0xe4, 0x30, 0x1d, 0xa8, 0xc1, 0x28, 0xe4, 0xc1, 0xc5, 0x0e,
	0x75, 0x90, 0x90, 0x3c, 0x86, 0x0e, 0x54, 0xa7, 0x4a, 0xe1, 0xb1, 0xce, 0x80, 0xd1, 0xc9, 0xe5,
	0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e,
	0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0xb4, 0xd2, 0x33, 0x4b, 0x32, 0x4a,
	0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x53, 0x52, 0x8b, 0x73, 0xf3, 0x8b, 0x75, 0x73, 0x12, 0x93,
	0x8a, 0xf5, 0xa1, 0x31, 0x54, 0x94, 0x5f, 0x5a, 0x92, 0x0a, 0x8d, 0xa8, 0x24, 0x36, 0x70, 0xf4,
	0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x09, 0xe4, 0x5d, 0xbf, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FilesServiceClient is the client API for FilesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FilesServiceClient interface {
	// UploadFile allows to upload a file
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (FilesService_UploadFileClient, error)
	// GetFile allows to get the contents of a file
	GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (FilesService_GetFileClient, error)
}

type filesServiceClient struct {
	cc grpc1.ClientConn
}

func NewFilesServiceClient(cc grpc1.ClientConn) FilesServiceClient {
	return &filesServiceClient{cc}
}

func (c *filesServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (FilesService_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_FilesService_serviceDesc.Streams[0], "/caerus.files.v1.FilesService/UploadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &filesServiceUploadFileClient{stream}
	return x, nil
}

type FilesService_UploadFileClient interface {
	Send(*FileChunk) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type filesServiceUploadFileClient struct {
	grpc.ClientStream
}

func (x *filesServiceUploadFileClient) Send(m *FileChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filesServiceUploadFileClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filesServiceClient) GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (FilesService_GetFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_FilesService_serviceDesc.Streams[1], "/caerus.files.v1.FilesService/GetFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &filesServiceGetFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FilesService_GetFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type filesServiceGetFileClient struct {
	grpc.ClientStream
}

func (x *filesServiceGetFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FilesServiceServer is the server API for FilesService service.
type FilesServiceServer interface {
	// UploadFile allows to upload a file
	UploadFile(FilesService_UploadFileServer) error
	// GetFile allows to get the contents of a file
	GetFile(*GetFileRequest, FilesService_GetFileServer) error
}

// UnimplementedFilesServiceServer can be embedded to have forward compatible implementations.
type UnimplementedFilesServiceServer struct {
}

func (*UnimplementedFilesServiceServer) UploadFile(srv FilesService_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (*UnimplementedFilesServiceServer) GetFile(req *GetFileRequest, srv FilesService_GetFileServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}

func RegisterFilesServiceServer(s grpc1.Server, srv FilesServiceServer) {
	s.RegisterService(&_FilesService_serviceDesc, srv)
}

func _FilesService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilesServiceServer).UploadFile(&filesServiceUploadFileServer{stream})
}

type FilesService_UploadFileServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*FileChunk, error)
	grpc.ServerStream
}

type filesServiceUploadFileServer struct {
	grpc.ServerStream
}

func (x *filesServiceUploadFileServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filesServiceUploadFileServer) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FilesService_GetFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FilesServiceServer).GetFile(m, &filesServiceGetFileServer{stream})
}

type FilesService_GetFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type filesServiceGetFileServer struct {
	grpc.ServerStream
}

func (x *filesServiceGetFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

var _FilesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "caerus.files.v1.FilesService",
	HandlerType: (*FilesServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FilesService_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetFile",
			Handler:       _FilesService_GetFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "caerus/files/v1/service.proto",
}

func (m *FileChunk) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FileChunk) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FileChunk) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintService(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.FileName) > 0 {
		i -= len(m.FileName)
		copy(dAtA[i:], m.FileName)
		i = encodeVarintService(dAtA, i, uint64(len(m.FileName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UploadFileResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UploadFileResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UploadFileResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Url) > 0 {
		i -= len(m.Url)
		copy(dAtA[i:], m.Url)
		i = encodeVarintService(dAtA, i, uint64(len(m.Url)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GetFileRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetFileRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetFileRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FileName) > 0 {
		i -= len(m.FileName)
		copy(dAtA[i:], m.FileName)
		i = encodeVarintService(dAtA, i, uint64(len(m.FileName)))
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
func (m *FileChunk) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FileName)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	return n
}

func (m *UploadFileResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Url)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	return n
}

func (m *GetFileRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FileName)
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
func (m *FileChunk) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: FileChunk: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FileChunk: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileName", wireType)
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
			m.FileName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
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
func (m *UploadFileResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UploadFileResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UploadFileResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Url", wireType)
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
			m.Url = string(dAtA[iNdEx:postIndex])
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
func (m *GetFileRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GetFileRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetFileRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FileName", wireType)
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
			m.FileName = string(dAtA[iNdEx:postIndex])
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