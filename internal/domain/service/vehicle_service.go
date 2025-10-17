package service

import (
	"context"
	"fmt"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
	"github.com/encomos/api-encomos/customer-service/internal/port/repository"
)

// VehicleService provides business logic for vehicle operations
type VehicleService struct {
	vehicleRepo  repository.VehicleRepository
	customerRepo repository.CustomerRepository
}

// NewVehicleService creates a new vehicle service
func NewVehicleService(
	vehicleRepo repository.VehicleRepository,
	customerRepo repository.CustomerRepository,
) *VehicleService {
	return &VehicleService{
		vehicleRepo:  vehicleRepo,
		customerRepo: customerRepo,
	}
}

// CreateVehicle creates a new vehicle with validation
func (s *VehicleService) CreateVehicle(ctx context.Context, create model.VehicleCreate) (*model.Vehicle, error) {
	// Verificar que el cliente existe
	_, err := s.customerRepo.GetByID(ctx, create.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Crear el vehículo
	vehicle := model.NewVehicle(create)
	if err := vehicle.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validar VIN si está presente
	if err := vehicle.ValidateVIN(); err != nil {
		return nil, fmt.Errorf("VIN validation error: %w", err)
	}

	// Verificar unicidad de VIN si está presente
	if vehicle.VIN != nil && *vehicle.VIN != "" {
		exists, err := s.vehicleRepo.ExistsByVIN(ctx, *vehicle.VIN, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check VIN uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("vehicle with VIN %s already exists", *vehicle.VIN)
		}
	}

	// Verificar unicidad de placa si está presente
	if vehicle.LicensePlate != nil && *vehicle.LicensePlate != "" {
		exists, err := s.vehicleRepo.ExistsByLicensePlate(ctx, *vehicle.LicensePlate, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check license plate uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("vehicle with license plate %s already exists", *vehicle.LicensePlate)
		}
	}

	// Crear el vehículo
	if err := s.vehicleRepo.Create(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	return vehicle, nil
}

// GetVehicle retrieves a vehicle by ID
func (s *VehicleService) GetVehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return vehicle, nil
}

// UpdateVehicle updates an existing vehicle
func (s *VehicleService) UpdateVehicle(ctx context.Context, update model.VehicleUpdate) (*model.Vehicle, error) {
	// Obtener el vehículo actual
	vehicle, err := s.vehicleRepo.GetByID(ctx, update.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle for update: %w", err)
	}

	// Verificar unicidad de VIN si se está cambiando
	if update.VIN != nil && *update.VIN != "" {
		if vehicle.VIN == nil || *vehicle.VIN != *update.VIN {
			exists, err := s.vehicleRepo.ExistsByVIN(ctx, *update.VIN, &update.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check VIN uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("vehicle with VIN %s already exists", *update.VIN)
			}
		}
	}

	// Verificar unicidad de placa si se está cambiando
	if update.LicensePlate != nil && *update.LicensePlate != "" {
		if vehicle.LicensePlate == nil || *vehicle.LicensePlate != *update.LicensePlate {
			exists, err := s.vehicleRepo.ExistsByLicensePlate(ctx, *update.LicensePlate, &update.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check license plate uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("vehicle with license plate %s already exists", *update.LicensePlate)
			}
		}
	}

	// Aplicar cambios
	vehicle.UpdateFromUpdate(update)

	// Validar después de los cambios
	if err := vehicle.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validar VIN después de los cambios
	if err := vehicle.ValidateVIN(); err != nil {
		return nil, fmt.Errorf("VIN validation error: %w", err)
	}

	// Actualizar en la base de datos
	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	return vehicle, nil
}

// DeleteVehicle deletes a vehicle
func (s *VehicleService) DeleteVehicle(ctx context.Context, id string) error {
	// Verificar que el vehículo existe
	_, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get vehicle for deletion: %w", err)
	}

	// Eliminar el vehículo
	if err := s.vehicleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	return nil
}

// ListVehicles lists vehicles with filtering and pagination
func (s *VehicleService) ListVehicles(ctx context.Context, filter model.VehicleFilter) ([]*model.Vehicle, int, error) {
	vehicles, total, err := s.vehicleRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}

	return vehicles, total, nil
}

// ListVehiclesByCustomer lists all vehicles for a customer
func (s *VehicleService) ListVehiclesByCustomer(ctx context.Context, customerID string) ([]*model.Vehicle, error) {
	// Verificar que el cliente existe
	_, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	vehicles, err := s.vehicleRepo.ListByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles by customer: %w", err)
	}

	return vehicles, nil
}

// GetVehicleByVIN retrieves a vehicle by VIN
func (s *VehicleService) GetVehicleByVIN(ctx context.Context, vin string) (*model.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByVIN(ctx, vin)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by VIN: %w", err)
	}

	return vehicle, nil
}

// GetVehicleByLicensePlate retrieves a vehicle by license plate
func (s *VehicleService) GetVehicleByLicensePlate(ctx context.Context, licensePlate string) (*model.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByLicensePlate(ctx, licensePlate)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by license plate: %w", err)
	}

	return vehicle, nil
}

// SearchVehicles searches vehicles by make, model, and year
func (s *VehicleService) SearchVehicles(ctx context.Context, make, model string, year *int) ([]*model.Vehicle, error) {
	vehicles, err := s.vehicleRepo.SearchByMakeModel(ctx, make, model, year)
	if err != nil {
		return nil, fmt.Errorf("failed to search vehicles: %w", err)
	}

	return vehicles, nil
}

// FindCompatibleVehicles finds vehicles compatible for parts (AutoParts feature)
func (s *VehicleService) FindCompatibleVehicles(ctx context.Context, make, model string, year int, yearRange int) ([]*model.Vehicle, error) {
	yearFrom := year - yearRange
	yearTo := year + yearRange

	vehicles, err := s.vehicleRepo.FindCompatibleVehicles(ctx, make, model, yearFrom, yearTo)
	if err != nil {
		return nil, fmt.Errorf("failed to find compatible vehicles: %w", err)
	}

	return vehicles, nil
}

// GetVehicleCompatibilityInfo returns compatibility information for a vehicle
func (s *VehicleService) GetVehicleCompatibilityInfo(ctx context.Context, id string) (map[string]interface{}, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	// Buscar vehículos compatibles (mismo make/model, años cercanos)
	compatibleVehicles, err := s.FindCompatibleVehicles(ctx, vehicle.Make, vehicle.Model, vehicle.Year, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to find compatible vehicles: %w", err)
	}

	// Buscar vehículos exactos (mismo make/model/año)
	exactMatches, err := s.vehicleRepo.ListByMakeModelYear(ctx, vehicle.Make, vehicle.Model, vehicle.Year)
	if err != nil {
		return nil, fmt.Errorf("failed to find exact matches: %w", err)
	}

	info := map[string]interface{}{
		"vehicle":              vehicle,
		"compatibility_string": vehicle.GetCompatibilityString(),
		"compatible_vehicles":  len(compatibleVehicles),
		"exact_matches":        len(exactMatches),
		"year_range":           fmt.Sprintf("%d-%d", vehicle.Year-3, vehicle.Year+3),
	}

	return info, nil
}

// ActivateVehicle activates a vehicle
func (s *VehicleService) ActivateVehicle(ctx context.Context, id string) error {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get vehicle: %w", err)
	}

	vehicle.Activate()

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to activate vehicle: %w", err)
	}

	return nil
}

// DeactivateVehicle deactivates a vehicle
func (s *VehicleService) DeactivateVehicle(ctx context.Context, id string) error {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get vehicle: %w", err)
	}

	vehicle.Deactivate()

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to deactivate vehicle: %w", err)
	}

	return nil
}

// CreateVehiclesForCustomer creates multiple vehicles for a customer in a batch
func (s *VehicleService) CreateVehiclesForCustomer(ctx context.Context, customerID string, vehicleCreates []model.VehicleCreate) ([]*model.Vehicle, error) {
	// Verificar que el cliente existe
	_, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	var vehicles []*model.Vehicle
	for _, create := range vehicleCreates {
		create.CustomerID = customerID // Asegurar que el customer ID esté configurado

		vehicle := model.NewVehicle(create)
		if err := vehicle.Validate(); err != nil {
			return nil, fmt.Errorf("validation error for vehicle %s %s: %w", vehicle.Make, vehicle.Model, err)
		}

		if err := vehicle.ValidateVIN(); err != nil {
			return nil, fmt.Errorf("VIN validation error for vehicle %s %s: %w", vehicle.Make, vehicle.Model, err)
		}

		vehicles = append(vehicles, vehicle)
	}

	// Verificar unicidad de VINs y placas antes de crear
	for _, vehicle := range vehicles {
		if vehicle.VIN != nil && *vehicle.VIN != "" {
			exists, err := s.vehicleRepo.ExistsByVIN(ctx, *vehicle.VIN, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to check VIN uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("vehicle with VIN %s already exists", *vehicle.VIN)
			}
		}

		if vehicle.LicensePlate != nil && *vehicle.LicensePlate != "" {
			exists, err := s.vehicleRepo.ExistsByLicensePlate(ctx, *vehicle.LicensePlate, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to check license plate uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("vehicle with license plate %s already exists", *vehicle.LicensePlate)
			}
		}
	}

	// Crear todos los vehículos en una transacción
	if err := s.vehicleRepo.CreateBatch(ctx, vehicles); err != nil {
		return nil, fmt.Errorf("failed to create vehicles batch: %w", err)
	}

	return vehicles, nil
}

// GetVehicleStats retrieves statistics for vehicles
func (s *VehicleService) GetVehicleStats(ctx context.Context) (map[string]interface{}, error) {
	totalVehicles, err := s.vehicleRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total vehicles: %w", err)
	}

	activeVehicles, err := s.vehicleRepo.CountActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count active vehicles: %w", err)
	}

	stats := map[string]interface{}{
		"total_vehicles":    totalVehicles,
		"active_vehicles":   activeVehicles,
		"inactive_vehicles": totalVehicles - activeVehicles,
	}

	return stats, nil
}
