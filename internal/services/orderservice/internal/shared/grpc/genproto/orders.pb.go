// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: orders.proto

package orders_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ShopItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=Title,proto3" json:"Title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	Quantity      uint64                 `protobuf:"varint,3,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=Price,proto3" json:"Price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShopItem) Reset() {
	*x = ShopItem{}
	mi := &file_orders_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShopItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShopItem) ProtoMessage() {}

func (x *ShopItem) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShopItem.ProtoReflect.Descriptor instead.
func (*ShopItem) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{0}
}

func (x *ShopItem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ShopItem) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ShopItem) GetQuantity() uint64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *ShopItem) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type Order struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	OrderID         string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	ShopItems       []*ShopItem            `protobuf:"bytes,2,rep,name=ShopItems,proto3" json:"ShopItems,omitempty"`
	Paid            bool                   `protobuf:"varint,3,opt,name=Paid,proto3" json:"Paid,omitempty"`
	Submitted       bool                   `protobuf:"varint,4,opt,name=Submitted,proto3" json:"Submitted,omitempty"`
	Completed       bool                   `protobuf:"varint,5,opt,name=Completed,proto3" json:"Completed,omitempty"`
	Canceled        bool                   `protobuf:"varint,6,opt,name=Canceled,proto3" json:"Canceled,omitempty"`
	TotalPrice      float64                `protobuf:"fixed64,7,opt,name=TotalPrice,proto3" json:"TotalPrice,omitempty"`
	AccountEmail    string                 `protobuf:"bytes,8,opt,name=AccountEmail,proto3" json:"AccountEmail,omitempty"`
	CancelReason    string                 `protobuf:"bytes,9,opt,name=CancelReason,proto3" json:"CancelReason,omitempty"`
	DeliveryAddress string                 `protobuf:"bytes,10,opt,name=DeliveryAddress,proto3" json:"DeliveryAddress,omitempty"`
	DeliveredTime   *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=DeliveredTime,proto3" json:"DeliveredTime,omitempty"`
	CreatedAt       *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt       *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	PaymentID       string                 `protobuf:"bytes,14,opt,name=PaymentID,proto3" json:"PaymentID,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_orders_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{1}
}

func (x *Order) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *Order) GetShopItems() []*ShopItem {
	if x != nil {
		return x.ShopItems
	}
	return nil
}

func (x *Order) GetPaid() bool {
	if x != nil {
		return x.Paid
	}
	return false
}

func (x *Order) GetSubmitted() bool {
	if x != nil {
		return x.Submitted
	}
	return false
}

func (x *Order) GetCompleted() bool {
	if x != nil {
		return x.Completed
	}
	return false
}

func (x *Order) GetCanceled() bool {
	if x != nil {
		return x.Canceled
	}
	return false
}

func (x *Order) GetTotalPrice() float64 {
	if x != nil {
		return x.TotalPrice
	}
	return 0
}

func (x *Order) GetAccountEmail() string {
	if x != nil {
		return x.AccountEmail
	}
	return ""
}

func (x *Order) GetCancelReason() string {
	if x != nil {
		return x.CancelReason
	}
	return ""
}

func (x *Order) GetDeliveryAddress() string {
	if x != nil {
		return x.DeliveryAddress
	}
	return ""
}

func (x *Order) GetDeliveredTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DeliveredTime
	}
	return nil
}

func (x *Order) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Order) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Order) GetPaymentID() string {
	if x != nil {
		return x.PaymentID
	}
	return ""
}

type OrderReadModel struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	ID              string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	OrderID         string                 `protobuf:"bytes,2,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	ShopItems       []*ShopItemReadModel   `protobuf:"bytes,3,rep,name=ShopItems,proto3" json:"ShopItems,omitempty"`
	Paid            bool                   `protobuf:"varint,4,opt,name=Paid,proto3" json:"Paid,omitempty"`
	Submitted       bool                   `protobuf:"varint,5,opt,name=Submitted,proto3" json:"Submitted,omitempty"`
	Completed       bool                   `protobuf:"varint,6,opt,name=Completed,proto3" json:"Completed,omitempty"`
	Canceled        bool                   `protobuf:"varint,7,opt,name=Canceled,proto3" json:"Canceled,omitempty"`
	TotalPrice      float64                `protobuf:"fixed64,8,opt,name=TotalPrice,proto3" json:"TotalPrice,omitempty"`
	AccountEmail    string                 `protobuf:"bytes,9,opt,name=AccountEmail,proto3" json:"AccountEmail,omitempty"`
	CancelReason    string                 `protobuf:"bytes,10,opt,name=CancelReason,proto3" json:"CancelReason,omitempty"`
	DeliveryAddress string                 `protobuf:"bytes,11,opt,name=DeliveryAddress,proto3" json:"DeliveryAddress,omitempty"`
	DeliveredTime   *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=DeliveredTime,proto3" json:"DeliveredTime,omitempty"`
	CreatedAt       *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt       *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	PaymentID       string                 `protobuf:"bytes,15,opt,name=PaymentID,proto3" json:"PaymentID,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *OrderReadModel) Reset() {
	*x = OrderReadModel{}
	mi := &file_orders_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderReadModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderReadModel) ProtoMessage() {}

func (x *OrderReadModel) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderReadModel.ProtoReflect.Descriptor instead.
func (*OrderReadModel) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{2}
}

func (x *OrderReadModel) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *OrderReadModel) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *OrderReadModel) GetShopItems() []*ShopItemReadModel {
	if x != nil {
		return x.ShopItems
	}
	return nil
}

func (x *OrderReadModel) GetPaid() bool {
	if x != nil {
		return x.Paid
	}
	return false
}

func (x *OrderReadModel) GetSubmitted() bool {
	if x != nil {
		return x.Submitted
	}
	return false
}

func (x *OrderReadModel) GetCompleted() bool {
	if x != nil {
		return x.Completed
	}
	return false
}

func (x *OrderReadModel) GetCanceled() bool {
	if x != nil {
		return x.Canceled
	}
	return false
}

func (x *OrderReadModel) GetTotalPrice() float64 {
	if x != nil {
		return x.TotalPrice
	}
	return 0
}

func (x *OrderReadModel) GetAccountEmail() string {
	if x != nil {
		return x.AccountEmail
	}
	return ""
}

func (x *OrderReadModel) GetCancelReason() string {
	if x != nil {
		return x.CancelReason
	}
	return ""
}

func (x *OrderReadModel) GetDeliveryAddress() string {
	if x != nil {
		return x.DeliveryAddress
	}
	return ""
}

func (x *OrderReadModel) GetDeliveredTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DeliveredTime
	}
	return nil
}

func (x *OrderReadModel) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *OrderReadModel) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *OrderReadModel) GetPaymentID() string {
	if x != nil {
		return x.PaymentID
	}
	return ""
}

type ShopItemReadModel struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=Title,proto3" json:"Title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	Quantity      uint64                 `protobuf:"varint,3,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=Price,proto3" json:"Price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShopItemReadModel) Reset() {
	*x = ShopItemReadModel{}
	mi := &file_orders_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShopItemReadModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShopItemReadModel) ProtoMessage() {}

func (x *ShopItemReadModel) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShopItemReadModel.ProtoReflect.Descriptor instead.
func (*ShopItemReadModel) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{3}
}

func (x *ShopItemReadModel) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ShopItemReadModel) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ShopItemReadModel) GetQuantity() uint64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *ShopItemReadModel) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type CreateOrderReq struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	AccountEmail    string                 `protobuf:"bytes,1,opt,name=AccountEmail,proto3" json:"AccountEmail,omitempty"`
	ShopItems       []*ShopItem            `protobuf:"bytes,2,rep,name=ShopItems,proto3" json:"ShopItems,omitempty"`
	DeliveryAddress string                 `protobuf:"bytes,3,opt,name=DeliveryAddress,proto3" json:"DeliveryAddress,omitempty"`
	DeliveryTime    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=DeliveryTime,proto3" json:"DeliveryTime,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *CreateOrderReq) Reset() {
	*x = CreateOrderReq{}
	mi := &file_orders_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderReq) ProtoMessage() {}

func (x *CreateOrderReq) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderReq.ProtoReflect.Descriptor instead.
func (*CreateOrderReq) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{4}
}

func (x *CreateOrderReq) GetAccountEmail() string {
	if x != nil {
		return x.AccountEmail
	}
	return ""
}

func (x *CreateOrderReq) GetShopItems() []*ShopItem {
	if x != nil {
		return x.ShopItems
	}
	return nil
}

func (x *CreateOrderReq) GetDeliveryAddress() string {
	if x != nil {
		return x.DeliveryAddress
	}
	return ""
}

func (x *CreateOrderReq) GetDeliveryTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DeliveryTime
	}
	return nil
}

type CreateOrderRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrderRes) Reset() {
	*x = CreateOrderRes{}
	mi := &file_orders_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRes) ProtoMessage() {}

func (x *CreateOrderRes) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRes.ProtoReflect.Descriptor instead.
func (*CreateOrderRes) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{5}
}

func (x *CreateOrderRes) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

type SubmitOrderReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitOrderReq) Reset() {
	*x = SubmitOrderReq{}
	mi := &file_orders_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitOrderReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitOrderReq) ProtoMessage() {}

func (x *SubmitOrderReq) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitOrderReq.ProtoReflect.Descriptor instead.
func (*SubmitOrderReq) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{6}
}

func (x *SubmitOrderReq) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

type SubmitOrderRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitOrderRes) Reset() {
	*x = SubmitOrderRes{}
	mi := &file_orders_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitOrderRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitOrderRes) ProtoMessage() {}

func (x *SubmitOrderRes) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitOrderRes.ProtoReflect.Descriptor instead.
func (*SubmitOrderRes) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{7}
}

func (x *SubmitOrderRes) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

type GetOrderByIDReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrderByIDReq) Reset() {
	*x = GetOrderByIDReq{}
	mi := &file_orders_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderByIDReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderByIDReq) ProtoMessage() {}

func (x *GetOrderByIDReq) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderByIDReq.ProtoReflect.Descriptor instead.
func (*GetOrderByIDReq) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{8}
}

func (x *GetOrderByIDReq) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type GetOrderByIDRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Order         *OrderReadModel        `protobuf:"bytes,1,opt,name=Order,proto3" json:"Order,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrderByIDRes) Reset() {
	*x = GetOrderByIDRes{}
	mi := &file_orders_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderByIDRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderByIDRes) ProtoMessage() {}

func (x *GetOrderByIDRes) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderByIDRes.ProtoReflect.Descriptor instead.
func (*GetOrderByIDRes) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{9}
}

func (x *GetOrderByIDRes) GetOrder() *OrderReadModel {
	if x != nil {
		return x.Order
	}
	return nil
}

type UpdateShoppingCartReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       string                 `protobuf:"bytes,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	ShopItems     []*ShopItem            `protobuf:"bytes,2,rep,name=ShopItems,proto3" json:"ShopItems,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateShoppingCartReq) Reset() {
	*x = UpdateShoppingCartReq{}
	mi := &file_orders_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateShoppingCartReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateShoppingCartReq) ProtoMessage() {}

func (x *UpdateShoppingCartReq) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateShoppingCartReq.ProtoReflect.Descriptor instead.
func (*UpdateShoppingCartReq) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateShoppingCartReq) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *UpdateShoppingCartReq) GetShopItems() []*ShopItem {
	if x != nil {
		return x.ShopItems
	}
	return nil
}

type UpdateShoppingCartRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateShoppingCartRes) Reset() {
	*x = UpdateShoppingCartRes{}
	mi := &file_orders_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateShoppingCartRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateShoppingCartRes) ProtoMessage() {}

func (x *UpdateShoppingCartRes) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateShoppingCartRes.ProtoReflect.Descriptor instead.
func (*UpdateShoppingCartRes) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{11}
}

type GetOrdersReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SearchText    string                 `protobuf:"bytes,1,opt,name=SearchText,proto3" json:"SearchText,omitempty"`
	Page          int32                  `protobuf:"varint,2,opt,name=Page,proto3" json:"Page,omitempty"`
	Size          int32                  `protobuf:"varint,3,opt,name=Size,proto3" json:"Size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrdersReq) Reset() {
	*x = GetOrdersReq{}
	mi := &file_orders_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrdersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersReq) ProtoMessage() {}

func (x *GetOrdersReq) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersReq.ProtoReflect.Descriptor instead.
func (*GetOrdersReq) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{12}
}

func (x *GetOrdersReq) GetSearchText() string {
	if x != nil {
		return x.SearchText
	}
	return ""
}

func (x *GetOrdersReq) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetOrdersReq) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

type GetOrdersRes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pagination    *Pagination            `protobuf:"bytes,1,opt,name=Pagination,proto3" json:"Pagination,omitempty"`
	Orders        []*OrderReadModel      `protobuf:"bytes,2,rep,name=Orders,proto3" json:"Orders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrdersRes) Reset() {
	*x = GetOrdersRes{}
	mi := &file_orders_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrdersRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersRes) ProtoMessage() {}

func (x *GetOrdersRes) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersRes.ProtoReflect.Descriptor instead.
func (*GetOrdersRes) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{13}
}

func (x *GetOrdersRes) GetPagination() *Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *GetOrdersRes) GetOrders() []*OrderReadModel {
	if x != nil {
		return x.Orders
	}
	return nil
}

type Pagination struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TotalItems    int64                  `protobuf:"varint,1,opt,name=TotalItems,proto3" json:"TotalItems,omitempty"`
	TotalPages    int32                  `protobuf:"varint,2,opt,name=TotalPages,proto3" json:"TotalPages,omitempty"`
	Page          int32                  `protobuf:"varint,3,opt,name=Page,proto3" json:"Page,omitempty"`
	Size          int32                  `protobuf:"varint,4,opt,name=Size,proto3" json:"Size,omitempty"`
	HasMore       bool                   `protobuf:"varint,5,opt,name=HasMore,proto3" json:"HasMore,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Pagination) Reset() {
	*x = Pagination{}
	mi := &file_orders_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Pagination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pagination) ProtoMessage() {}

func (x *Pagination) ProtoReflect() protoreflect.Message {
	mi := &file_orders_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pagination.ProtoReflect.Descriptor instead.
func (*Pagination) Descriptor() ([]byte, []int) {
	return file_orders_proto_rawDescGZIP(), []int{14}
}

func (x *Pagination) GetTotalItems() int64 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *Pagination) GetTotalPages() int32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

func (x *Pagination) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *Pagination) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Pagination) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

var File_orders_proto protoreflect.FileDescriptor

const file_orders_proto_rawDesc = "" +
	"\n" +
	"\forders.proto\x12\x0eorders_service\x1a\x1fgoogle/protobuf/timestamp.proto\"t\n" +
	"\bShopItem\x12\x14\n" +
	"\x05Title\x18\x01 \x01(\tR\x05Title\x12 \n" +
	"\vDescription\x18\x02 \x01(\tR\vDescription\x12\x1a\n" +
	"\bQuantity\x18\x03 \x01(\x04R\bQuantity\x12\x14\n" +
	"\x05Price\x18\x04 \x01(\x01R\x05Price\"\xab\x04\n" +
	"\x05Order\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\tR\aOrderID\x126\n" +
	"\tShopItems\x18\x02 \x03(\v2\x18.orders_service.ShopItemR\tShopItems\x12\x12\n" +
	"\x04Paid\x18\x03 \x01(\bR\x04Paid\x12\x1c\n" +
	"\tSubmitted\x18\x04 \x01(\bR\tSubmitted\x12\x1c\n" +
	"\tCompleted\x18\x05 \x01(\bR\tCompleted\x12\x1a\n" +
	"\bCanceled\x18\x06 \x01(\bR\bCanceled\x12\x1e\n" +
	"\n" +
	"TotalPrice\x18\a \x01(\x01R\n" +
	"TotalPrice\x12\"\n" +
	"\fAccountEmail\x18\b \x01(\tR\fAccountEmail\x12\"\n" +
	"\fCancelReason\x18\t \x01(\tR\fCancelReason\x12(\n" +
	"\x0fDeliveryAddress\x18\n" +
	" \x01(\tR\x0fDeliveryAddress\x12@\n" +
	"\rDeliveredTime\x18\v \x01(\v2\x1a.google.protobuf.TimestampR\rDeliveredTime\x128\n" +
	"\tCreatedAt\x18\f \x01(\v2\x1a.google.protobuf.TimestampR\tCreatedAt\x128\n" +
	"\tUpdatedAt\x18\r \x01(\v2\x1a.google.protobuf.TimestampR\tUpdatedAt\x12\x1c\n" +
	"\tPaymentID\x18\x0e \x01(\tR\tPaymentID\"\xcd\x04\n" +
	"\x0eOrderReadModel\x12\x0e\n" +
	"\x02ID\x18\x01 \x01(\tR\x02ID\x12\x18\n" +
	"\aOrderID\x18\x02 \x01(\tR\aOrderID\x12?\n" +
	"\tShopItems\x18\x03 \x03(\v2!.orders_service.ShopItemReadModelR\tShopItems\x12\x12\n" +
	"\x04Paid\x18\x04 \x01(\bR\x04Paid\x12\x1c\n" +
	"\tSubmitted\x18\x05 \x01(\bR\tSubmitted\x12\x1c\n" +
	"\tCompleted\x18\x06 \x01(\bR\tCompleted\x12\x1a\n" +
	"\bCanceled\x18\a \x01(\bR\bCanceled\x12\x1e\n" +
	"\n" +
	"TotalPrice\x18\b \x01(\x01R\n" +
	"TotalPrice\x12\"\n" +
	"\fAccountEmail\x18\t \x01(\tR\fAccountEmail\x12\"\n" +
	"\fCancelReason\x18\n" +
	" \x01(\tR\fCancelReason\x12(\n" +
	"\x0fDeliveryAddress\x18\v \x01(\tR\x0fDeliveryAddress\x12@\n" +
	"\rDeliveredTime\x18\f \x01(\v2\x1a.google.protobuf.TimestampR\rDeliveredTime\x128\n" +
	"\tCreatedAt\x18\r \x01(\v2\x1a.google.protobuf.TimestampR\tCreatedAt\x128\n" +
	"\tUpdatedAt\x18\x0e \x01(\v2\x1a.google.protobuf.TimestampR\tUpdatedAt\x12\x1c\n" +
	"\tPaymentID\x18\x0f \x01(\tR\tPaymentID\"}\n" +
	"\x11ShopItemReadModel\x12\x14\n" +
	"\x05Title\x18\x01 \x01(\tR\x05Title\x12 \n" +
	"\vDescription\x18\x02 \x01(\tR\vDescription\x12\x1a\n" +
	"\bQuantity\x18\x03 \x01(\x04R\bQuantity\x12\x14\n" +
	"\x05Price\x18\x04 \x01(\x01R\x05Price\"\xd6\x01\n" +
	"\x0eCreateOrderReq\x12\"\n" +
	"\fAccountEmail\x18\x01 \x01(\tR\fAccountEmail\x126\n" +
	"\tShopItems\x18\x02 \x03(\v2\x18.orders_service.ShopItemR\tShopItems\x12(\n" +
	"\x0fDeliveryAddress\x18\x03 \x01(\tR\x0fDeliveryAddress\x12>\n" +
	"\fDeliveryTime\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\fDeliveryTime\"*\n" +
	"\x0eCreateOrderRes\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\tR\aOrderID\"*\n" +
	"\x0eSubmitOrderReq\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\tR\aOrderID\"*\n" +
	"\x0eSubmitOrderRes\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\tR\aOrderID\"!\n" +
	"\x0fGetOrderByIDReq\x12\x0e\n" +
	"\x02ID\x18\x01 \x01(\tR\x02ID\"G\n" +
	"\x0fGetOrderByIDRes\x124\n" +
	"\x05Order\x18\x01 \x01(\v2\x1e.orders_service.OrderReadModelR\x05Order\"i\n" +
	"\x15UpdateShoppingCartReq\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\tR\aOrderID\x126\n" +
	"\tShopItems\x18\x02 \x03(\v2\x18.orders_service.ShopItemR\tShopItems\"\x17\n" +
	"\x15UpdateShoppingCartRes\"V\n" +
	"\fGetOrdersReq\x12\x1e\n" +
	"\n" +
	"SearchText\x18\x01 \x01(\tR\n" +
	"SearchText\x12\x12\n" +
	"\x04Page\x18\x02 \x01(\x05R\x04Page\x12\x12\n" +
	"\x04Size\x18\x03 \x01(\x05R\x04Size\"\x82\x01\n" +
	"\fGetOrdersRes\x12:\n" +
	"\n" +
	"Pagination\x18\x01 \x01(\v2\x1a.orders_service.PaginationR\n" +
	"Pagination\x126\n" +
	"\x06Orders\x18\x02 \x03(\v2\x1e.orders_service.OrderReadModelR\x06Orders\"\x8e\x01\n" +
	"\n" +
	"Pagination\x12\x1e\n" +
	"\n" +
	"TotalItems\x18\x01 \x01(\x03R\n" +
	"TotalItems\x12\x1e\n" +
	"\n" +
	"TotalPages\x18\x02 \x01(\x05R\n" +
	"TotalPages\x12\x12\n" +
	"\x04Page\x18\x03 \x01(\x05R\x04Page\x12\x12\n" +
	"\x04Size\x18\x04 \x01(\x05R\x04Size\x12\x18\n" +
	"\aHasMore\x18\x05 \x01(\bR\aHasMore2\xac\x03\n" +
	"\rOrdersService\x12M\n" +
	"\vCreateOrder\x12\x1e.orders_service.CreateOrderReq\x1a\x1e.orders_service.CreateOrderRes\x12M\n" +
	"\vSubmitOrder\x12\x1e.orders_service.SubmitOrderReq\x1a\x1e.orders_service.SubmitOrderRes\x12b\n" +
	"\x12UpdateShoppingCart\x12%.orders_service.UpdateShoppingCartReq\x1a%.orders_service.UpdateShoppingCartRes\x12P\n" +
	"\fGetOrderByID\x12\x1f.orders_service.GetOrderByIDReq\x1a\x1f.orders_service.GetOrderByIDRes\x12G\n" +
	"\tGetOrders\x12\x1c.orders_service.GetOrdersReq\x1a\x1c.orders_service.GetOrdersResB\x13Z\x11./;orders_serviceb\x06proto3"

var (
	file_orders_proto_rawDescOnce sync.Once
	file_orders_proto_rawDescData []byte
)

func file_orders_proto_rawDescGZIP() []byte {
	file_orders_proto_rawDescOnce.Do(func() {
		file_orders_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_orders_proto_rawDesc), len(file_orders_proto_rawDesc)))
	})
	return file_orders_proto_rawDescData
}

var file_orders_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_orders_proto_goTypes = []any{
	(*ShopItem)(nil),              // 0: orders_service.ShopItem
	(*Order)(nil),                 // 1: orders_service.Order
	(*OrderReadModel)(nil),        // 2: orders_service.OrderReadModel
	(*ShopItemReadModel)(nil),     // 3: orders_service.ShopItemReadModel
	(*CreateOrderReq)(nil),        // 4: orders_service.CreateOrderReq
	(*CreateOrderRes)(nil),        // 5: orders_service.CreateOrderRes
	(*SubmitOrderReq)(nil),        // 6: orders_service.SubmitOrderReq
	(*SubmitOrderRes)(nil),        // 7: orders_service.SubmitOrderRes
	(*GetOrderByIDReq)(nil),       // 8: orders_service.GetOrderByIDReq
	(*GetOrderByIDRes)(nil),       // 9: orders_service.GetOrderByIDRes
	(*UpdateShoppingCartReq)(nil), // 10: orders_service.UpdateShoppingCartReq
	(*UpdateShoppingCartRes)(nil), // 11: orders_service.UpdateShoppingCartRes
	(*GetOrdersReq)(nil),          // 12: orders_service.GetOrdersReq
	(*GetOrdersRes)(nil),          // 13: orders_service.GetOrdersRes
	(*Pagination)(nil),            // 14: orders_service.Pagination
	(*timestamppb.Timestamp)(nil), // 15: google.protobuf.Timestamp
}
var file_orders_proto_depIdxs = []int32{
	0,  // 0: orders_service.Order.ShopItems:type_name -> orders_service.ShopItem
	15, // 1: orders_service.Order.DeliveredTime:type_name -> google.protobuf.Timestamp
	15, // 2: orders_service.Order.CreatedAt:type_name -> google.protobuf.Timestamp
	15, // 3: orders_service.Order.UpdatedAt:type_name -> google.protobuf.Timestamp
	3,  // 4: orders_service.OrderReadModel.ShopItems:type_name -> orders_service.ShopItemReadModel
	15, // 5: orders_service.OrderReadModel.DeliveredTime:type_name -> google.protobuf.Timestamp
	15, // 6: orders_service.OrderReadModel.CreatedAt:type_name -> google.protobuf.Timestamp
	15, // 7: orders_service.OrderReadModel.UpdatedAt:type_name -> google.protobuf.Timestamp
	0,  // 8: orders_service.CreateOrderReq.ShopItems:type_name -> orders_service.ShopItem
	15, // 9: orders_service.CreateOrderReq.DeliveryTime:type_name -> google.protobuf.Timestamp
	2,  // 10: orders_service.GetOrderByIDRes.Order:type_name -> orders_service.OrderReadModel
	0,  // 11: orders_service.UpdateShoppingCartReq.ShopItems:type_name -> orders_service.ShopItem
	14, // 12: orders_service.GetOrdersRes.Pagination:type_name -> orders_service.Pagination
	2,  // 13: orders_service.GetOrdersRes.Orders:type_name -> orders_service.OrderReadModel
	4,  // 14: orders_service.OrdersService.CreateOrder:input_type -> orders_service.CreateOrderReq
	6,  // 15: orders_service.OrdersService.SubmitOrder:input_type -> orders_service.SubmitOrderReq
	10, // 16: orders_service.OrdersService.UpdateShoppingCart:input_type -> orders_service.UpdateShoppingCartReq
	8,  // 17: orders_service.OrdersService.GetOrderByID:input_type -> orders_service.GetOrderByIDReq
	12, // 18: orders_service.OrdersService.GetOrders:input_type -> orders_service.GetOrdersReq
	5,  // 19: orders_service.OrdersService.CreateOrder:output_type -> orders_service.CreateOrderRes
	7,  // 20: orders_service.OrdersService.SubmitOrder:output_type -> orders_service.SubmitOrderRes
	11, // 21: orders_service.OrdersService.UpdateShoppingCart:output_type -> orders_service.UpdateShoppingCartRes
	9,  // 22: orders_service.OrdersService.GetOrderByID:output_type -> orders_service.GetOrderByIDRes
	13, // 23: orders_service.OrdersService.GetOrders:output_type -> orders_service.GetOrdersRes
	19, // [19:24] is the sub-list for method output_type
	14, // [14:19] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_orders_proto_init() }
func file_orders_proto_init() {
	if File_orders_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_orders_proto_rawDesc), len(file_orders_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orders_proto_goTypes,
		DependencyIndexes: file_orders_proto_depIdxs,
		MessageInfos:      file_orders_proto_msgTypes,
	}.Build()
	File_orders_proto = out.File
	file_orders_proto_goTypes = nil
	file_orders_proto_depIdxs = nil
}
