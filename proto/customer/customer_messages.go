// Additional request/response messages for customer.proto
// This file extends the basic protobuf messages with all required types

package customerpb

import (
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// CreateCustomerRequest for creating a new customer
type CreateCustomerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName    string                 `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string                 `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Email        string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phone        string                 `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	CustomerType string                 `protobuf:"bytes,5,opt,name=customer_type,json=customerType,proto3" json:"customer_type,omitempty"`
	CompanyName  string                 `protobuf:"bytes,6,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	TaxId        string                 `protobuf:"bytes,7,opt,name=tax_id,json=taxId,proto3" json:"tax_id,omitempty"`
	Address      string                 `protobuf:"bytes,8,opt,name=address,proto3" json:"address,omitempty"`
	Birthday     *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Notes        string                 `protobuf:"bytes,10,opt,name=notes,proto3" json:"notes,omitempty"`
	Preferences  *structpb.Struct       `protobuf:"bytes,11,opt,name=preferences,proto3" json:"preferences,omitempty"`
	Vehicles     []*CreateVehicleRequest `protobuf:"bytes,12,rep,name=vehicles,proto3" json:"vehicles,omitempty"`
}

func (x *CreateCustomerRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *CreateCustomerRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *CreateCustomerRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateCustomerRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *CreateCustomerRequest) GetCustomerType() string {
	if x != nil {
		return x.CustomerType
	}
	return ""
}

func (x *CreateCustomerRequest) GetCompanyName() string {
	if x != nil {
		return x.CompanyName
	}
	return ""
}

func (x *CreateCustomerRequest) GetTaxId() string {
	if x != nil {
		return x.TaxId
	}
	return ""
}

func (x *CreateCustomerRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *CreateCustomerRequest) GetBirthday() *timestamppb.Timestamp {
	if x != nil {
		return x.Birthday
	}
	return nil
}

func (x *CreateCustomerRequest) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *CreateCustomerRequest) GetPreferences() *structpb.Struct {
	if x != nil {
		return x.Preferences
	}
	return nil
}

func (x *CreateCustomerRequest) GetVehicles() []*CreateVehicleRequest {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

// CreateCustomerResponse for customer creation response
type CreateCustomerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Customer *Customer `protobuf:"bytes,1,opt,name=customer,proto3" json:"customer,omitempty"`
}

func (x *CreateCustomerResponse) GetCustomer() *Customer {
	if x != nil {
		return x.Customer
	}
	return nil
}

// UpdateCustomerRequest for updating a customer
type UpdateCustomerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName    string                 `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string                 `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Email        string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Phone        string                 `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	CustomerType string                 `protobuf:"bytes,6,opt,name=customer_type,json=customerType,proto3" json:"customer_type,omitempty"`
	CompanyName  string                 `protobuf:"bytes,7,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	TaxId        string                 `protobuf:"bytes,8,opt,name=tax_id,json=taxId,proto3" json:"tax_id,omitempty"`
	Address      string                 `protobuf:"bytes,9,opt,name=address,proto3" json:"address,omitempty"`
	Birthday     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Notes        string                 `protobuf:"bytes,11,opt,name=notes,proto3" json:"notes,omitempty"`
	Preferences  *structpb.Struct       `protobuf:"bytes,12,opt,name=preferences,proto3" json:"preferences,omitempty"`
	IsActive     bool                   `protobuf:"varint,13,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

func (x *UpdateCustomerRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCustomerRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UpdateCustomerRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UpdateCustomerRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UpdateCustomerRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *UpdateCustomerRequest) GetCustomerType() string {
	if x != nil {
		return x.CustomerType
	}
	return ""
}

func (x *UpdateCustomerRequest) GetCompanyName() string {
	if x != nil {
		return x.CompanyName
	}
	return ""
}

func (x *UpdateCustomerRequest) GetTaxId() string {
	if x != nil {
		return x.TaxId
	}
	return ""
}

func (x *UpdateCustomerRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *UpdateCustomerRequest) GetBirthday() *timestamppb.Timestamp {
	if x != nil {
		return x.Birthday
	}
	return nil
}

func (x *UpdateCustomerRequest) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *UpdateCustomerRequest) GetPreferences() *structpb.Struct {
	if x != nil {
		return x.Preferences
	}
	return nil
}

func (x *UpdateCustomerRequest) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

// UpdateCustomerResponse for customer update response
type UpdateCustomerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Customer *Customer `protobuf:"bytes,1,opt,name=customer,proto3" json:"customer,omitempty"`
}

func (x *UpdateCustomerResponse) GetCustomer() *Customer {
	if x != nil {
		return x.Customer
	}
	return nil
}

// DeleteCustomerRequest for deleting a customer
type DeleteCustomerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteCustomerRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// DeleteCustomerResponse for customer deletion response
type DeleteCustomerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *DeleteCustomerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// Vehicle-related messages
type ListVehiclesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId int64  `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Search     string `protobuf:"bytes,2,opt,name=search,proto3" json:"search,omitempty"`
	ActiveOnly bool   `protobuf:"varint,3,opt,name=active_only,json=activeOnly,proto3" json:"active_only,omitempty"`
	Page       int32  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	Limit      int32  `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListVehiclesRequest) GetCustomerId() int64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *ListVehiclesRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *ListVehiclesRequest) GetActiveOnly() bool {
	if x != nil {
		return x.ActiveOnly
	}
	return false
}

func (x *ListVehiclesRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListVehiclesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ListVehiclesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicles []*Vehicle `protobuf:"bytes,1,rep,name=vehicles,proto3" json:"vehicles,omitempty"`
	Total    int32      `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *ListVehiclesResponse) GetVehicles() []*Vehicle {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

func (x *ListVehiclesResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

// Additional message types for Vehicle operations
type GetVehicleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetVehicleRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetVehicleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicle *Vehicle `protobuf:"bytes,1,opt,name=vehicle,proto3" json:"vehicle,omitempty"`
}

func (x *GetVehicleResponse) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

type CreateVehicleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId   int64            `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Make         string           `protobuf:"bytes,2,opt,name=make,proto3" json:"make,omitempty"`
	Model        string           `protobuf:"bytes,3,opt,name=model,proto3" json:"model,omitempty"`
	Year         int32            `protobuf:"varint,4,opt,name=year,proto3" json:"year,omitempty"`
	Vin          string           `protobuf:"bytes,5,opt,name=vin,proto3" json:"vin,omitempty"`
	LicensePlate string           `protobuf:"bytes,6,opt,name=license_plate,json=licensePlate,proto3" json:"license_plate,omitempty"`
	Color        string           `protobuf:"bytes,7,opt,name=color,proto3" json:"color,omitempty"`
	Engine       string           `protobuf:"bytes,8,opt,name=engine,proto3" json:"engine,omitempty"`
	Notes        string           `protobuf:"bytes,9,opt,name=notes,proto3" json:"notes,omitempty"`
	Metadata     *structpb.Struct `protobuf:"bytes,10,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *CreateVehicleRequest) GetCustomerId() int64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *CreateVehicleRequest) GetMake() string {
	if x != nil {
		return x.Make
	}
	return ""
}

func (x *CreateVehicleRequest) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *CreateVehicleRequest) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *CreateVehicleRequest) GetVin() string {
	if x != nil {
		return x.Vin
	}
	return ""
}

func (x *CreateVehicleRequest) GetLicensePlate() string {
	if x != nil {
		return x.LicensePlate
	}
	return ""
}

func (x *CreateVehicleRequest) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *CreateVehicleRequest) GetEngine() string {
	if x != nil {
		return x.Engine
	}
	return ""
}

func (x *CreateVehicleRequest) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *CreateVehicleRequest) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type CreateVehicleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicle *Vehicle `protobuf:"bytes,1,opt,name=vehicle,proto3" json:"vehicle,omitempty"`
}

func (x *CreateVehicleResponse) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

// Search-related messages
type SearchCustomersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query        string `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	SearchFields string `protobuf:"bytes,2,opt,name=search_fields,json=searchFields,proto3" json:"search_fields,omitempty"`
	Limit        int32  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *SearchCustomersRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchCustomersRequest) GetSearchFields() string {
	if x != nil {
		return x.SearchFields
	}
	return ""
}

func (x *SearchCustomersRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type SearchCustomersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Customers []*Customer `protobuf:"bytes,1,rep,name=customers,proto3" json:"customers,omitempty"`
	Total     int32       `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *SearchCustomersResponse) GetCustomers() []*Customer {
	if x != nil {
		return x.Customers
	}
	return nil
}

func (x *SearchCustomersResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

// Customer notes and history
type AddCustomerNoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId int64  `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Note       string `protobuf:"bytes,2,opt,name=note,proto3" json:"note,omitempty"`
	Type       string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *AddCustomerNoteRequest) GetCustomerId() int64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *AddCustomerNoteRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *AddCustomerNoteRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type AddCustomerNoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Note *CustomerNote `protobuf:"bytes,1,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *AddCustomerNoteResponse) GetNote() *CustomerNote {
	if x != nil {
		return x.Note
	}
	return nil
}

// Customer History (placeholder for future integration)
type GetCustomerHistoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId int64                  `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type       string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	DateFrom   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=date_from,json=dateFrom,proto3" json:"date_from,omitempty"`
	DateTo     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date_to,json=dateTo,proto3" json:"date_to,omitempty"`
	Page       int32                  `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`
	Limit      int32                  `protobuf:"varint,6,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetCustomerHistoryRequest) GetCustomerId() int64 {
	if x != nil {
		return x.CustomerId
	}
	return 0
}

func (x *GetCustomerHistoryRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *GetCustomerHistoryRequest) GetDateFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.DateFrom
	}
	return nil
}

func (x *GetCustomerHistoryRequest) GetDateTo() *timestamppb.Timestamp {
	if x != nil {
		return x.DateTo
	}
	return nil
}

func (x *GetCustomerHistoryRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetCustomerHistoryRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type CustomerHistoryItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type        string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Title       string                 `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Description string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Amount      float64                `protobuf:"fixed64,5,opt,name=amount,proto3" json:"amount,omitempty"`
	Status      string                 `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	Data        *structpb.Struct       `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *CustomerHistoryItem) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CustomerHistoryItem) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *CustomerHistoryItem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CustomerHistoryItem) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CustomerHistoryItem) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *CustomerHistoryItem) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CustomerHistoryItem) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *CustomerHistoryItem) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type GetCustomerHistoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*CustomerHistoryItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Total int32                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *GetCustomerHistoryResponse) GetItems() []*CustomerHistoryItem {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *GetCustomerHistoryResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

// Additional Vehicle request/response messages
type UpdateVehicleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int64            `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Make         string           `protobuf:"bytes,2,opt,name=make,proto3" json:"make,omitempty"`
	Model        string           `protobuf:"bytes,3,opt,name=model,proto3" json:"model,omitempty"`
	Year         int32            `protobuf:"varint,4,opt,name=year,proto3" json:"year,omitempty"`
	Vin          string           `protobuf:"bytes,5,opt,name=vin,proto3" json:"vin,omitempty"`
	LicensePlate string           `protobuf:"bytes,6,opt,name=license_plate,json=licensePlate,proto3" json:"license_plate,omitempty"`
	Color        string           `protobuf:"bytes,7,opt,name=color,proto3" json:"color,omitempty"`
	Engine       string           `protobuf:"bytes,8,opt,name=engine,proto3" json:"engine,omitempty"`
	Notes        string           `protobuf:"bytes,9,opt,name=notes,proto3" json:"notes,omitempty"`
	IsActive     bool             `protobuf:"varint,10,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	Metadata     *structpb.Struct `protobuf:"bytes,11,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *UpdateVehicleRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateVehicleRequest) GetMake() string {
	if x != nil {
		return x.Make
	}
	return ""
}

func (x *UpdateVehicleRequest) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *UpdateVehicleRequest) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *UpdateVehicleRequest) GetVin() string {
	if x != nil {
		return x.Vin
	}
	return ""
}

func (x *UpdateVehicleRequest) GetLicensePlate() string {
	if x != nil {
		return x.LicensePlate
	}
	return ""
}

func (x *UpdateVehicleRequest) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *UpdateVehicleRequest) GetEngine() string {
	if x != nil {
		return x.Engine
	}
	return ""
}

func (x *UpdateVehicleRequest) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *UpdateVehicleRequest) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

func (x *UpdateVehicleRequest) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type UpdateVehicleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicle *Vehicle `protobuf:"bytes,1,opt,name=vehicle,proto3" json:"vehicle,omitempty"`
}

func (x *UpdateVehicleResponse) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

type DeleteVehicleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteVehicleRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteVehicleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *DeleteVehicleResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}
