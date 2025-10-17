package service

import (
	"context"
	"fmt"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
	"github.com/encomos/api-encomos/customer-service/internal/port/repository"
)

// CustomerService provides business logic for customer operations
type CustomerService struct {
	customerRepo     repository.CustomerRepository
	vehicleRepo      repository.VehicleRepository
	customerNoteRepo repository.CustomerNoteRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(
	customerRepo repository.CustomerRepository,
	vehicleRepo repository.VehicleRepository,
	customerNoteRepo repository.CustomerNoteRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo:     customerRepo,
		vehicleRepo:      vehicleRepo,
		customerNoteRepo: customerNoteRepo,
	}
}

// CreateCustomer creates a new customer with validation
func (s *CustomerService) CreateCustomer(ctx context.Context, create model.CustomerCreate) (*model.Customer, error) {
	// Validar datos de entrada
	customer := model.NewCustomer(create)
	if err := customer.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Verificar unicidad de email si está presente
	if customer.Email != nil && *customer.Email != "" {
		exists, err := s.customerRepo.ExistsByEmail(ctx, *customer.Email, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("customer with email %s already exists", *customer.Email)
		}
	}

	// Verificar unicidad de Tax ID si está presente
	if customer.TaxID != nil && *customer.TaxID != "" {
		exists, err := s.customerRepo.ExistsByTaxID(ctx, *customer.TaxID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check tax ID uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("customer with tax ID %s already exists", *customer.TaxID)
		}
	}

	// Crear el cliente
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetCustomer retrieves a customer by ID with optional related data
func (s *CustomerService) GetCustomer(ctx context.Context, id string, includeVehicles, includeNotes bool) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Cargar vehículos si se solicita
	if includeVehicles {
		vehicles, err := s.vehicleRepo.ListByCustomer(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to load customer vehicles: %w", err)
		}
		customer.Vehicles = vehicles
	}

	// Cargar notas si se solicita
	if includeNotes {
		notes, err := s.customerNoteRepo.ListRecentByCustomer(ctx, id, 10) // Últimas 10 notas
		if err != nil {
			return nil, fmt.Errorf("failed to load customer notes: %w", err)
		}
		customer.CustomerNotes = notes
	}

	return customer, nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(ctx context.Context, update model.CustomerUpdate) (*model.Customer, error) {
	// Obtener el cliente actual
	customer, err := s.customerRepo.GetByID(ctx, update.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer for update: %w", err)
	}

	// Verificar unicidad de email si se está cambiando
	if update.Email != nil && *update.Email != "" {
		if customer.Email == nil || *customer.Email != *update.Email {
			exists, err := s.customerRepo.ExistsByEmail(ctx, *update.Email, &update.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("customer with email %s already exists", *update.Email)
			}
		}
	}

	// Verificar unicidad de Tax ID si se está cambiando
	if update.TaxID != nil && *update.TaxID != "" {
		if customer.TaxID == nil || *customer.TaxID != *update.TaxID {
			exists, err := s.customerRepo.ExistsByTaxID(ctx, *update.TaxID, &update.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to check tax ID uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("customer with tax ID %s already exists", *update.TaxID)
			}
		}
	}

	// Aplicar cambios
	customer.UpdateFromUpdate(update)

	// Validar después de los cambios
	if err := customer.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Actualizar en la base de datos
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

// DeleteCustomer deletes a customer (soft delete by deactivating)
func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	// Verificar que el cliente existe
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer for deletion: %w", err)
	}

	// Verificar si el cliente tiene vehículos activos
	vehicles, err := s.vehicleRepo.ListActiveByCustomer(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check customer vehicles: %w", err)
	}

	if len(vehicles) > 0 {
		// Soft delete - desactivar en lugar de eliminar
		customer.Deactivate()
		if err := s.customerRepo.Update(ctx, customer); err != nil {
			return fmt.Errorf("failed to deactivate customer: %w", err)
		}
		return nil
	}

	// Hard delete si no tiene vehículos
	if err := s.customerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// ListCustomers lists customers with filtering and pagination
func (s *CustomerService) ListCustomers(ctx context.Context, filter model.CustomerFilter) ([]*model.Customer, int, error) {
	customers, total, err := s.customerRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}

	return customers, total, nil
}

// SearchCustomers performs advanced search on customers
func (s *CustomerService) SearchCustomers(ctx context.Context, filter model.CustomerSearchFilter) ([]*model.Customer, error) {
	if filter.Query == "" {
		return []*model.Customer{}, nil
	}

	customers, err := s.customerRepo.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}

	return customers, nil
}

// GetCustomerByEmail retrieves a customer by email
func (s *CustomerService) GetCustomerByEmail(ctx context.Context, email string) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by email: %w", err)
	}

	return customer, nil
}

// GetCustomerByTaxID retrieves a customer by tax ID
func (s *CustomerService) GetCustomerByTaxID(ctx context.Context, taxID string) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByTaxID(ctx, taxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by tax ID: %w", err)
	}

	return customer, nil
}

// ActivateCustomer activates a customer
func (s *CustomerService) ActivateCustomer(ctx context.Context, id string) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	customer.Activate()

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return fmt.Errorf("failed to activate customer: %w", err)
	}

	return nil
}

// DeactivateCustomer deactivates a customer
func (s *CustomerService) DeactivateCustomer(ctx context.Context, id string) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	customer.Deactivate()

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return fmt.Errorf("failed to deactivate customer: %w", err)
	}

	return nil
}

// GetCustomerStats retrieves statistics for customers
func (s *CustomerService) GetCustomerStats(ctx context.Context) (map[string]interface{}, error) {
	totalCustomers, err := s.customerRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total customers: %w", err)
	}

	activeCustomers, err := s.customerRepo.CountActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count active customers: %w", err)
	}

	individualCustomers, err := s.customerRepo.CountByType(ctx, model.CustomerTypeIndividual)
	if err != nil {
		return nil, fmt.Errorf("failed to count individual customers: %w", err)
	}

	businessCustomers, err := s.customerRepo.CountByType(ctx, model.CustomerTypeBusiness)
	if err != nil {
		return nil, fmt.Errorf("failed to count business customers: %w", err)
	}

	stats := map[string]interface{}{
		"total_customers":      totalCustomers,
		"active_customers":     activeCustomers,
		"inactive_customers":   totalCustomers - activeCustomers,
		"individual_customers": individualCustomers,
		"business_customers":   businessCustomers,
		"calculated_at":        time.Now(),
	}

	return stats, nil
}

// AddCustomerNote adds a note to a customer
func (s *CustomerService) AddCustomerNote(ctx context.Context, create model.CustomerNoteCreate) (*model.CustomerNote, error) {
	// Verificar que el cliente existe
	_, err := s.customerRepo.GetByID(ctx, create.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Crear la nota
	note := model.NewCustomerNote(create)
	if err := note.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := s.customerNoteRepo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("failed to create customer note: %w", err)
	}

	return note, nil
}

// GetCustomerNotes retrieves notes for a customer
func (s *CustomerService) GetCustomerNotes(ctx context.Context, customerID string, noteType string, limit int) ([]*model.CustomerNote, error) {
	// Verificar que el cliente existe
	_, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	if noteType != "" {
		return s.customerNoteRepo.ListByCustomerAndType(ctx, customerID, noteType)
	}

	if limit > 0 {
		return s.customerNoteRepo.ListRecentByCustomer(ctx, customerID, limit)
	}

	return s.customerNoteRepo.ListByCustomer(ctx, customerID)
}

// SetCustomerPreference sets a preference for a customer
func (s *CustomerService) SetCustomerPreference(ctx context.Context, customerID string, key string, value interface{}) error {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	customer.SetPreference(key, value)

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return fmt.Errorf("failed to update customer preference: %w", err)
	}

	return nil
}

// GetCustomerPreference gets a preference for a customer
func (s *CustomerService) GetCustomerPreference(ctx context.Context, customerID string, key string) (interface{}, error) {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	value, exists := customer.GetPreference(key)
	if !exists {
		return nil, fmt.Errorf("preference %s not found for customer", key)
	}

	return value, nil
}
