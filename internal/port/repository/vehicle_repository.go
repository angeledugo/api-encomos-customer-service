package repository

import (
	"context"

	"github.com/yourorg/api-encomos/customer-service/internal/domain/model"
)

// VehicleRepository define la interfaz para operaciones de repositorio de vehículos
type VehicleRepository interface {
	// CRUD básico
	Create(ctx context.Context, vehicle *model.Vehicle) error
	GetByID(ctx context.Context, id int64) (*model.Vehicle, error)
	Update(ctx context.Context, vehicle *model.Vehicle) error
	Delete(ctx context.Context, id int64) error

	// Búsquedas
	List(ctx context.Context, filter model.VehicleFilter) ([]*model.Vehicle, int, error)
	ListByCustomer(ctx context.Context, customerID int64) ([]*model.Vehicle, error)
	
	// Búsquedas específicas
	GetByVIN(ctx context.Context, vin string) (*model.Vehicle, error)
	GetByLicensePlate(ctx context.Context, licensePlate string) (*model.Vehicle, error)
	SearchByMakeModel(ctx context.Context, make, model string, year *int) ([]*model.Vehicle, error)
	
	// Consultas de compatibilidad (para AutoParts)
	FindCompatibleVehicles(ctx context.Context, make, model string, yearFrom, yearTo int) ([]*model.Vehicle, error)
	ListByMakeModelYear(ctx context.Context, make, model string, year int) ([]*model.Vehicle, error)
	
	// Estadísticas
	Count(ctx context.Context) (int64, error)
	CountByCustomer(ctx context.Context, customerID int64) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	
	// Validaciones
	ExistsByVIN(ctx context.Context, vin string, excludeID *int64) (bool, error)
	ExistsByLicensePlate(ctx context.Context, licensePlate string, excludeID *int64) (bool, error)
	
	// Operaciones en lote
	CreateBatch(ctx context.Context, vehicles []*model.Vehicle) error
	ListActiveByCustomer(ctx context.Context, customerID int64) ([]*model.Vehicle, error)
}
