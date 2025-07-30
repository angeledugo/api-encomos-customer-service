package repository

import (
	"context"

	"github.com/yourorg/api-encomos/customer-service/internal/domain/model"
)

// CustomerRepository define la interfaz para operaciones de repositorio de clientes
type CustomerRepository interface {
	// CRUD básico
	Create(ctx context.Context, customer *model.Customer) error
	GetByID(ctx context.Context, id int64) (*model.Customer, error)
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, id int64) error

	// Búsquedas
	List(ctx context.Context, filter model.CustomerFilter) ([]*model.Customer, int, error)
	Search(ctx context.Context, filter model.CustomerSearchFilter) ([]*model.Customer, error)
	GetByEmail(ctx context.Context, email string) (*model.Customer, error)
	GetByTaxID(ctx context.Context, taxID string) (*model.Customer, error)
	
	// Consultas específicas
	ListByType(ctx context.Context, customerType string, page, limit int) ([]*model.Customer, int, error)
	ListActive(ctx context.Context, page, limit int) ([]*model.Customer, int, error)
	ListInactive(ctx context.Context, page, limit int) ([]*model.Customer, int, error)
	
	// Estadísticas
	Count(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, customerType string) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	
	// Validaciones
	ExistsByEmail(ctx context.Context, email string, excludeID *int64) (bool, error)
	ExistsByTaxID(ctx context.Context, taxID string, excludeID *int64) (bool, error)
}
