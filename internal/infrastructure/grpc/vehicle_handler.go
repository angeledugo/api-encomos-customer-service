package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yourorg/api-encomos/customer-service/internal/domain/model"
	"github.com/yourorg/api-encomos/customer-service/internal/domain/service"
	customerpb "github.com/yourorg/api-encomos/customer-service/proto/customer"
)

// VehicleHandler handles vehicle-related gRPC requests
type VehicleHandler struct {
	vehicleService *service.VehicleService
}

// NewVehicleHandler creates a new vehicle handler
func NewVehicleHandler(vehicleService *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
	}
}

// ListVehicles lists vehicles with filtering and pagination
func (h *VehicleHandler) ListVehicles(ctx context.Context, req *customerpb.ListVehiclesRequest) (*customerpb.ListVehiclesResponse, error) {
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
	filter := model.VehicleFilter{
		CustomerID: req.CustomerId,
		Search:     req.Search,
		ActiveOnly: req.ActiveOnly,
		Page:       int(req.Page),
		Limit:      int(req.Limit),
	}

	// Ejecutar búsqueda
	vehicles, total, err := h.vehicleService.ListVehicles(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vehicles: %v", err)
	}

	// Convertir a protobuf
	pbVehicles := make([]*customerpb.Vehicle, len(vehicles))
	for i, vehicle := range vehicles {
		pbVehicles[i] = h.vehicleToProto(vehicle)
	}

	return &customerpb.ListVehiclesResponse{
		Vehicles: pbVehicles,
		Total:    int32(total),
	}, nil
}

// GetVehicle retrieves a vehicle by ID
func (h *VehicleHandler) GetVehicle(ctx context.Context, req *customerpb.GetVehicleRequest) (*customerpb.GetVehicleResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "vehicle ID must be positive")
	}

	vehicle, err := h.vehicleService.GetVehicle(ctx, req.Id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "vehicle not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get vehicle: %v", err)
	}

	return &customerpb.GetVehicleResponse{
		Vehicle: h.vehicleToProto(vehicle),
	}, nil
}

// CreateVehicle creates a new vehicle
func (h *VehicleHandler) CreateVehicle(ctx context.Context, req *customerpb.CreateVehicleRequest) (*customerpb.CreateVehicleResponse, error) {
	// Validar entrada
	if req.CustomerId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "customer ID must be positive")
	}
	if req.Make == "" {
		return nil, status.Errorf(codes.InvalidArgument, "make is required")
	}
	if req.Model == "" {
		return nil, status.Errorf(codes.InvalidArgument, "model is required")
	}
	if req.Year <= 1900 || req.Year > 2100 {
		return nil, status.Errorf(codes.InvalidArgument, "year must be between 1900 and 2100")
	}

	// Convertir de protobuf a modelo
	create := model.VehicleCreate{
		CustomerID:   req.CustomerId,
		Make:         req.Make,
		Model:        req.Model,
		Year:         int(req.Year),
		VIN:          stringPtrFromProto(req.Vin),
		LicensePlate: stringPtrFromProto(req.LicensePlate),
		Color:        stringPtrFromProto(req.Color),
		Engine:       stringPtrFromProto(req.Engine),
		Notes:        stringPtrFromProto(req.Notes),
		Metadata:     make(model.VehicleMetadata),
	}

	if req.Metadata != nil {
		create.Metadata = req.Metadata.AsMap()
	}

	// Crear vehículo
	vehicle, err := h.vehicleService.CreateVehicle(ctx, create)
	if err != nil {
		if isValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		if isDuplicateError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "vehicle already exists: %v", err)
		}
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create vehicle: %v", err)
	}

	return &customerpb.CreateVehicleResponse{
		Vehicle: h.vehicleToProto(vehicle),
	}, nil
}

// UpdateVehicle updates an existing vehicle
func (h *VehicleHandler) UpdateVehicle(ctx context.Context, req *customerpb.UpdateVehicleRequest) (*customerpb.UpdateVehicleResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "vehicle ID must be positive")
	}

	// Convertir de protobuf a modelo
	update := model.VehicleUpdate{
		ID: req.Id,
	}

	if req.Make != "" {
		update.Make = &req.Make
	}
	if req.Model != "" {
		update.Model = &req.Model
	}
	if req.Year > 0 {
		year := int(req.Year)
		update.Year = &year
	}
	if req.Vin != "" {
		update.VIN = &req.Vin
	}
	if req.LicensePlate != "" {
		update.LicensePlate = &req.LicensePlate
	}
	if req.Color != "" {
		update.Color = &req.Color
	}
	if req.Engine != "" {
		update.Engine = &req.Engine
	}
	if req.Notes != "" {
		update.Notes = &req.Notes
	}

	update.IsActive = &req.IsActive

	if req.Metadata != nil {
		metadata := req.Metadata.AsMap()
		update.Metadata = metadata
	}

	// Actualizar vehículo
	vehicle, err := h.vehicleService.UpdateVehicle(ctx, update)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "vehicle not found")
		}
		if isValidationError(err) {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		if isDuplicateError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "vehicle already exists: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update vehicle: %v", err)
	}

	return &customerpb.UpdateVehicleResponse{
		Vehicle: h.vehicleToProto(vehicle),
	}, nil
}

// DeleteVehicle deletes a vehicle
func (h *VehicleHandler) DeleteVehicle(ctx context.Context, req *customerpb.DeleteVehicleRequest) (*customerpb.DeleteVehicleResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "vehicle ID must be positive")
	}

	err := h.vehicleService.DeleteVehicle(ctx, req.Id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "vehicle not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete vehicle: %v", err)
	}

	return &customerpb.DeleteVehicleResponse{
		Success: true,
	}, nil
}

// vehicleToProto converts a domain Vehicle to protobuf
func (h *VehicleHandler) vehicleToProto(vehicle *model.Vehicle) *customerpb.Vehicle {
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

	// TODO: Convert metadata to protobuf Struct when needed

	return pb
}
