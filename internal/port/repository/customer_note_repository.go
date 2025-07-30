package repository

import (
	"context"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
)

// CustomerNoteRepository define la interfaz para operaciones de repositorio de notas de clientes
type CustomerNoteRepository interface {
	// CRUD básico
	Create(ctx context.Context, note *model.CustomerNote) error
	GetByID(ctx context.Context, id int64) (*model.CustomerNote, error)
	Delete(ctx context.Context, id int64) error

	// Búsquedas
	List(ctx context.Context, filter model.CustomerNoteFilter) ([]*model.CustomerNote, int, error)
	ListByCustomer(ctx context.Context, customerID int64) ([]*model.CustomerNote, error)
	ListByCustomerAndType(ctx context.Context, customerID int64, noteType string) ([]*model.CustomerNote, error)

	// Consultas específicas
	ListByStaff(ctx context.Context, staffID int64, page, limit int) ([]*model.CustomerNote, int, error)
	ListByType(ctx context.Context, noteType string, page, limit int) ([]*model.CustomerNote, int, error)
	ListRecent(ctx context.Context, limit int) ([]*model.CustomerNote, error)

	// Búsquedas por fecha
	ListByDateRange(ctx context.Context, customerID int64, from, to *time.Time) ([]*model.CustomerNote, error)
	ListRecentByCustomer(ctx context.Context, customerID int64, limit int) ([]*model.CustomerNote, error)

	// Estadísticas
	Count(ctx context.Context) (int64, error)
	CountByCustomer(ctx context.Context, customerID int64) (int64, error)
	CountByType(ctx context.Context, noteType string) (int64, error)
	CountByStaff(ctx context.Context, staffID int64) (int64, error)

	// Análisis
	GetNoteTypesCount(ctx context.Context, customerID int64) (map[string]int64, error)
	GetMostActiveStaff(ctx context.Context, limit int) ([]map[string]interface{}, error)
}
