package grpc

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
	"github.com/encomos/api-encomos/customer-service/internal/domain/service"
	customerpb "github.com/encomos/api-encomos/customer-service/proto/customer"
)

// CustomerHandler handles customer-related gRPC requests
type CustomerHandler struct {
	customerpb.UnimplementedCustomerServiceServer
	customerService *service.CustomerService
	vehicleService  *service.VehicleService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService *service.CustomerService, vehicleService *service.VehicleService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		vehicleService:  vehicleService,
	}
}

// ListCustomers lists customers with filtering and pagination
func (h *CustomerHandler) ListCustomers(ctx context.Context, req *customerpb.ListCustomersRequest) (*customerpb.ListCustomersResponse, error) {
	// Validar entrada
	if req.Page < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "page must be non-negative")
	}
	if req.Limit <= 0 {
		req.Limit = 20 // Default limit
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	// Construir filtro
	filter := model.CustomerFilter{
		Search:       req.Search,
		CustomerType: req.CustomerType,
		ActiveOnly:   req.ActiveOnly,
		Page:         int(req.Page),
		Limit:        int(req.Limit),
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
	}

	// Ejecutar búsqueda
	customers, total, err := h.customerService.ListCustomers(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list customers: %v", err)
	}

	// Convertir a protobuf
	pbCustomers := make([]*customerpb.Customer, len(customers))
	for i, customer := range customers {
		pbCustomers[i] = h.customerToProto(customer)
	}

	// Calcular páginas totales
	totalPages := int32((total + int(req.Limit) - 1) / int(req.Limit))

	return &customerpb.ListCustomersResponse{
		Customers:  pbCustomers,
		Total:      int32(total),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetCustomer retrieves a customer by ID
func (h *CustomerHandler) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}

	customer, err := h.customerService.GetCustomer(ctx, req.Id, req.IncludeVehicles, req.IncludeNotes)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get customer: %v", err)
	}

	return &customerpb.GetCustomerResponse{
		Customer: h.customerToProto(customer),
	}, nil
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	// Validar entrada
	if req.FirstName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "first name is required")
	}
	if req.LastName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "last name is required")
	}
	if req.CustomerType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer type is required")
	}

	// Extraer tenant ID del contexto
	tenantID, err := extractTenantIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extract tenant ID: %v", err)
	}

	// Convertir de protobuf a modelo
	create := model.CustomerCreate{
		TenantID:     fmt.Sprintf("%d", tenantID),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        stringPtrFromProto(req.Email),
		Phone:        stringPtrFromProto(req.Phone),
		CustomerType: req.CustomerType,
		CompanyName:  stringPtrFromProto(req.CompanyName),
		TaxID:        stringPtrFromProto(req.TaxId),
		Address:      stringPtrFromProto(req.Address),
		Notes:        stringPtrFromProto(req.Notes),
		Preferences:  make(model.CustomerPreferences),
	}

	if req.Birthday != nil {
		birthday := req.Birthday.AsTime()
		create.Birthday = &birthday
	}

	if req.Preferences != nil {
		create.Preferences = req.Preferences.AsMap()
	}

	// Crear cliente
	customer, err := h.customerService.CreateCustomer(ctx, create)
	if err != nil {
		if isValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		if isDuplicateError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "customer already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create customer: %v", err)
	}

	return &customerpb.CreateCustomerResponse{
		Customer: h.customerToProto(customer),
	}, nil
}

// UpdateCustomer updates an existing customer
func (h *CustomerHandler) UpdateCustomer(ctx context.Context, req *customerpb.UpdateCustomerRequest) (*customerpb.UpdateCustomerResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}

	// Convertir de protobuf a modelo
	update := model.CustomerUpdate{
		ID: req.Id,
	}

	if req.FirstName != "" {
		update.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		update.LastName = &req.LastName
	}
	if req.Email != "" {
		update.Email = &req.Email
	}
	if req.Phone != "" {
		update.Phone = &req.Phone
	}
	if req.CustomerType != "" {
		update.CustomerType = &req.CustomerType
	}
	if req.CompanyName != "" {
		update.CompanyName = &req.CompanyName
	}
	if req.TaxId != "" {
		update.TaxID = &req.TaxId
	}
	if req.Address != "" {
		update.Address = &req.Address
	}
	if req.Notes != "" {
		update.Notes = &req.Notes
	}

	update.IsActive = &req.IsActive

	if req.Birthday != nil {
		birthday := req.Birthday.AsTime()
		update.Birthday = &birthday
	}

	if req.Preferences != nil {
		prefs := req.Preferences.AsMap()
		update.Preferences = prefs
	}

	// Actualizar cliente
	customer, err := h.customerService.UpdateCustomer(ctx, update)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		if isValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		if isDuplicateError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "customer already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update customer: %v", err)
	}

	return &customerpb.UpdateCustomerResponse{
		Customer: h.customerToProto(customer),
	}, nil
}

// DeleteCustomer deletes a customer
func (h *CustomerHandler) DeleteCustomer(ctx context.Context, req *customerpb.DeleteCustomerRequest) (*customerpb.DeleteCustomerResponse, error) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}

	err := h.customerService.DeleteCustomer(ctx, req.Id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete customer: %v", err)
	}

	return &customerpb.DeleteCustomerResponse{
		Success: true,
	}, nil
}

// SearchCustomers performs advanced search on customers
func (h *CustomerHandler) SearchCustomers(ctx context.Context, req *customerpb.SearchCustomersRequest) (*customerpb.SearchCustomersResponse, error) {
	if req.Query == "" {
		return &customerpb.SearchCustomersResponse{
			Customers: []*customerpb.Customer{},
			Total:     0,
		}, nil
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	// Construir filtro de búsqueda
	var searchFields []string
	if req.SearchFields != "" {
		// Parsear campos de búsqueda separados por coma
		// Por simplicidad, usaremos todos los campos por defecto
		searchFields = []string{"name", "email", "phone", "tax_id"}
	}

	filter := model.CustomerSearchFilter{
		Query:        req.Query,
		SearchFields: searchFields,
		Limit:        limit,
	}

	// Ejecutar búsqueda
	customers, err := h.customerService.SearchCustomers(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search customers: %v", err)
	}

	// Convertir a protobuf
	pbCustomers := make([]*customerpb.Customer, len(customers))
	for i, customer := range customers {
		pbCustomers[i] = h.customerToProto(customer)
	}

	return &customerpb.SearchCustomersResponse{
		Customers: pbCustomers,
		Total:     int32(len(customers)),
	}, nil
}

// AddCustomerNote adds a note to a customer
func (h *CustomerHandler) AddCustomerNote(ctx context.Context, req *customerpb.AddCustomerNoteRequest) (*customerpb.AddCustomerNoteResponse, error) {
	if req.CustomerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}
	if req.Note == "" {
		return nil, status.Errorf(codes.InvalidArgument, "note content is required")
	}

	// Por ahora, usar valores dummy para staff info - esto se obtendrá del token JWT en producción
	create := model.CustomerNoteCreate{
		CustomerID: req.CustomerId,
		StaffID:    "1",           // TODO: Obtener del contexto de autenticación
		StaffName:  "System User", // TODO: Obtener del contexto de autenticación
		Note:       req.Note,
		Type:       req.Type,
	}

	if create.Type == "" {
		create.Type = "general"
	}

	note, err := h.customerService.AddCustomerNote(ctx, create)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		if isValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to add customer note: %v", err)
	}

	return &customerpb.AddCustomerNoteResponse{
		Note: h.customerNoteToProto(note),
	}, nil
}

// GetCustomerHistory retrieves customer history (placeholder implementation)
func (h *CustomerHandler) GetCustomerHistory(ctx context.Context, req *customerpb.GetCustomerHistoryRequest) (*customerpb.GetCustomerHistoryResponse, error) {
	if req.CustomerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID is required")
	}

	// TODO: Implementar lógica real de historial cuando tengamos integración con sales/appointments
	return &customerpb.GetCustomerHistoryResponse{
		Items: []*customerpb.CustomerHistoryItem{},
		Total: 0,
	}, nil
}

// customerToProto converts a domain Customer to protobuf
func (h *CustomerHandler) customerToProto(customer *model.Customer) *customerpb.Customer {
	pb := &customerpb.Customer{
		Id:           customer.ID,
		TenantId:     customer.TenantID,
		FirstName:    customer.FirstName,
		LastName:     customer.LastName,
		CustomerType: customer.CustomerType,
		IsActive:     customer.IsActive,
		CreatedAt:    timestamppb.New(customer.CreatedAt),
		UpdatedAt:    timestamppb.New(customer.UpdatedAt),
	}

	if customer.Email != nil {
		pb.Email = *customer.Email
	}
	if customer.Phone != nil {
		pb.Phone = *customer.Phone
	}
	if customer.CompanyName != nil {
		pb.CompanyName = *customer.CompanyName
	}
	if customer.TaxID != nil {
		pb.TaxId = *customer.TaxID
	}
	if customer.Address != nil {
		pb.Address = *customer.Address
	}
	if customer.Notes != nil {
		pb.Notes = *customer.Notes
	}
	if customer.Birthday != nil {
		pb.Birthday = timestamppb.New(*customer.Birthday)
	}

	// Convert preferences
	if customer.Preferences != nil && len(customer.Preferences) > 0 {
		// TODO: Convert map to protobuf Struct
	}

	// Convert vehicles if present
	if customer.Vehicles != nil {
		pb.Vehicles = make([]*customerpb.Vehicle, len(customer.Vehicles))
		for i, vehicle := range customer.Vehicles {
			pb.Vehicles[i] = h.vehicleToProto(vehicle)
		}
	}

	// Convert notes if present
	if customer.CustomerNotes != nil {
		pb.CustomerNotes = make([]*customerpb.CustomerNote, len(customer.CustomerNotes))
		for i, note := range customer.CustomerNotes {
			pb.CustomerNotes[i] = h.customerNoteToProto(note)
		}
	}

	return pb
}

// vehicleToProto converts a domain Vehicle to protobuf
func (h *CustomerHandler) vehicleToProto(vehicle *model.Vehicle) *customerpb.Vehicle {
	pb := &customerpb.Vehicle{
		Id:         vehicle.ID,
		CustomerId: vehicle.CustomerID,
		Make:       vehicle.Make,
		Model:      vehicle.Model,
		Year:       int32(vehicle.Year),
		IsActive:   vehicle.IsActive,
		CreatedAt:  timestamppb.New(vehicle.CreatedAt),
		UpdatedAt:  timestamppb.New(vehicle.UpdatedAt),
	}

	if vehicle.VIN != nil {
		pb.Vin = *vehicle.VIN
	}
	if vehicle.LicensePlate != nil {
		pb.LicensePlate = *vehicle.LicensePlate
	}
	if vehicle.Color != nil {
		pb.Color = *vehicle.Color
	}
	if vehicle.Engine != nil {
		pb.Engine = *vehicle.Engine
	}
	if vehicle.Notes != nil {
		pb.Notes = *vehicle.Notes
	}

	// TODO: Convert metadata to protobuf Struct

	return pb
}

// customerNoteToProto converts a domain CustomerNote to protobuf
func (h *CustomerHandler) customerNoteToProto(note *model.CustomerNote) *customerpb.CustomerNote {
	return &customerpb.CustomerNote{
		Id:         note.ID,
		CustomerId: note.CustomerID,
		StaffId:    note.StaffID,
		StaffName:  note.StaffName,
		Note:       note.Note,
		Type:       note.Type,
		CreatedAt:  timestamppb.New(note.CreatedAt),
	}
}

// Helper functions

func extractTenantIDFromContext(ctx context.Context) (int64, error) {
	// TODO: Implementar extracción real del tenant ID desde el contexto/JWT
	// Por ahora, retornar un valor dummy
	return 1, nil
}

func stringPtrFromProto(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func isNotFoundError(err error) bool {
	return err == sql.ErrNoRows ||
		(err != nil && (containsString(err.Error(), "not found") ||
			containsString(err.Error(), "does not exist")))
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*model.ValidationError)
	return ok || containsString(err.Error(), "validation error")
}

func isDuplicateError(err error) bool {
	return err != nil && (containsString(err.Error(), "already exists") ||
		containsString(err.Error(), "duplicate") ||
		containsString(err.Error(), "unique constraint"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RegisterService registers the customer service with the gRPC server
func (h *CustomerHandler) RegisterService(server *grpc.Server) {
	customerpb.RegisterCustomerServiceServer(server, h)
}
