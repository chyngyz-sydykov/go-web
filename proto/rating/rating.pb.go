// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: rating/rating.proto

package rating

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for saving a rating
type SaveRatingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookId  int32  `protobuf:"varint,1,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"` // ID of the book being rated
	Rating  int32  `protobuf:"varint,2,opt,name=rating,proto3" json:"rating,omitempty"`               // Rating from 1 to 5
	Comment string `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`              // Optional comment
}

func (x *SaveRatingRequest) Reset() {
	*x = SaveRatingRequest{}
	mi := &file_rating_rating_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveRatingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveRatingRequest) ProtoMessage() {}

func (x *SaveRatingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rating_rating_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveRatingRequest.ProtoReflect.Descriptor instead.
func (*SaveRatingRequest) Descriptor() ([]byte, []int) {
	return file_rating_rating_proto_rawDescGZIP(), []int{0}
}

func (x *SaveRatingRequest) GetBookId() int32 {
	if x != nil {
		return x.BookId
	}
	return 0
}

func (x *SaveRatingRequest) GetRating() int32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

func (x *SaveRatingRequest) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

// Response message for saving a rating, returning the saved rating details
type SaveRatingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rating *Rating `protobuf:"bytes,1,opt,name=rating,proto3" json:"rating,omitempty"` // The saved rating details
}

func (x *SaveRatingResponse) Reset() {
	*x = SaveRatingResponse{}
	mi := &file_rating_rating_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SaveRatingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveRatingResponse) ProtoMessage() {}

func (x *SaveRatingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rating_rating_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveRatingResponse.ProtoReflect.Descriptor instead.
func (*SaveRatingResponse) Descriptor() ([]byte, []int) {
	return file_rating_rating_proto_rawDescGZIP(), []int{1}
}

func (x *SaveRatingResponse) GetRating() *Rating {
	if x != nil {
		return x.Rating
	}
	return nil
}

// Request message for retrieving ratings of a specific book
type GetRatingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BookId int32 `protobuf:"varint,1,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"` // ID of the book for which ratings are requested
}

func (x *GetRatingsRequest) Reset() {
	*x = GetRatingsRequest{}
	mi := &file_rating_rating_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRatingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRatingsRequest) ProtoMessage() {}

func (x *GetRatingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rating_rating_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRatingsRequest.ProtoReflect.Descriptor instead.
func (*GetRatingsRequest) Descriptor() ([]byte, []int) {
	return file_rating_rating_proto_rawDescGZIP(), []int{2}
}

func (x *GetRatingsRequest) GetBookId() int32 {
	if x != nil {
		return x.BookId
	}
	return 0
}

// Response message containing a list of ratings
type GetRatingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ratings []*Rating `protobuf:"bytes,1,rep,name=ratings,proto3" json:"ratings,omitempty"` // List of ratings for the specified book
}

func (x *GetRatingsResponse) Reset() {
	*x = GetRatingsResponse{}
	mi := &file_rating_rating_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRatingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRatingsResponse) ProtoMessage() {}

func (x *GetRatingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rating_rating_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRatingsResponse.ProtoReflect.Descriptor instead.
func (*GetRatingsResponse) Descriptor() ([]byte, []int) {
	return file_rating_rating_proto_rawDescGZIP(), []int{3}
}

func (x *GetRatingsResponse) GetRatings() []*Rating {
	if x != nil {
		return x.Ratings
	}
	return nil
}

// Rating message that represents a single rating entry
type Rating struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RatingId string `protobuf:"bytes,1,opt,name=rating_id,json=ratingId,proto3" json:"rating_id,omitempty"` // Unique ID of the rating (UUID format)
	BookId   int32  `protobuf:"varint,2,opt,name=book_id,json=bookId,proto3" json:"book_id,omitempty"`      // ID of the book this rating is associated with
	Rating   int32  `protobuf:"varint,3,opt,name=rating,proto3" json:"rating,omitempty"`                    // Rating from 1 to 5
	Comment  string `protobuf:"bytes,4,opt,name=comment,proto3" json:"comment,omitempty"`                   // Optional comment
}

func (x *Rating) Reset() {
	*x = Rating{}
	mi := &file_rating_rating_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rating) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rating) ProtoMessage() {}

func (x *Rating) ProtoReflect() protoreflect.Message {
	mi := &file_rating_rating_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rating.ProtoReflect.Descriptor instead.
func (*Rating) Descriptor() ([]byte, []int) {
	return file_rating_rating_proto_rawDescGZIP(), []int{4}
}

func (x *Rating) GetRatingId() string {
	if x != nil {
		return x.RatingId
	}
	return ""
}

func (x *Rating) GetBookId() int32 {
	if x != nil {
		return x.BookId
	}
	return 0
}

func (x *Rating) GetRating() int32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

func (x *Rating) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

var File_rating_rating_proto protoreflect.FileDescriptor

var file_rating_rating_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x5e, 0x0a,
	0x11, 0x53, 0x61, 0x76, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x62, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3c, 0x0a,
	0x12, 0x53, 0x61, 0x76, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x52, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x2c, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x62, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x22, 0x3e, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x28, 0x0a, 0x07, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x07, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x70, 0x0a, 0x06, 0x52, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x62, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x32, 0x99, 0x01, 0x0a, 0x0d,
	0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a,
	0x0a, 0x53, 0x61, 0x76, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x2e, 0x72, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2e,
	0x53, 0x61, 0x76, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x19, 0x2e, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x61, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x72, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x79, 0x6e, 0x67, 0x79, 0x7a, 0x2d, 0x73, 0x79,
	0x64, 0x79, 0x6b, 0x6f, 0x76, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x2d, 0x72, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x3b,
	0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rating_rating_proto_rawDescOnce sync.Once
	file_rating_rating_proto_rawDescData = file_rating_rating_proto_rawDesc
)

func file_rating_rating_proto_rawDescGZIP() []byte {
	file_rating_rating_proto_rawDescOnce.Do(func() {
		file_rating_rating_proto_rawDescData = protoimpl.X.CompressGZIP(file_rating_rating_proto_rawDescData)
	})
	return file_rating_rating_proto_rawDescData
}

var file_rating_rating_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_rating_rating_proto_goTypes = []any{
	(*SaveRatingRequest)(nil),  // 0: rating.SaveRatingRequest
	(*SaveRatingResponse)(nil), // 1: rating.SaveRatingResponse
	(*GetRatingsRequest)(nil),  // 2: rating.GetRatingsRequest
	(*GetRatingsResponse)(nil), // 3: rating.GetRatingsResponse
	(*Rating)(nil),             // 4: rating.Rating
}
var file_rating_rating_proto_depIdxs = []int32{
	4, // 0: rating.SaveRatingResponse.rating:type_name -> rating.Rating
	4, // 1: rating.GetRatingsResponse.ratings:type_name -> rating.Rating
	0, // 2: rating.RatingService.SaveRating:input_type -> rating.SaveRatingRequest
	2, // 3: rating.RatingService.GetRatings:input_type -> rating.GetRatingsRequest
	1, // 4: rating.RatingService.SaveRating:output_type -> rating.SaveRatingResponse
	3, // 5: rating.RatingService.GetRatings:output_type -> rating.GetRatingsResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rating_rating_proto_init() }
func file_rating_rating_proto_init() {
	if File_rating_rating_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rating_rating_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rating_rating_proto_goTypes,
		DependencyIndexes: file_rating_rating_proto_depIdxs,
		MessageInfos:      file_rating_rating_proto_msgTypes,
	}.Build()
	File_rating_rating_proto = out.File
	file_rating_rating_proto_rawDesc = nil
	file_rating_rating_proto_goTypes = nil
	file_rating_rating_proto_depIdxs = nil
}
