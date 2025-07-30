package repository

import (
	"context"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
)

// CustomerStatsRepository define la interfaz para operaciones de repositorio de estadísticas de clientes
type CustomerStatsRepository interface {
	// CRUD básico
	Create(ctx context.Context, stats *model.CustomerStats) error
	GetByCustomerID(ctx context.Context, customerID int64) (*model.CustomerStats, error)
	Update(ctx context.Context, stats *model.CustomerStats) error
	Delete(ctx context.Context, customerID int64) error

	// Operaciones de cálculo
	CalculateAndSave(ctx context.Context, customerID int64) (*model.CustomerStats, error)
	RecalculateAll(ctx context.Context) error
	RecalculateOutdated(ctx context.Context) error

	// Consultas estadísticas
	ListTopCustomersBySpent(ctx context.Context, limit int) ([]*model.CustomerStats, error)
	ListTopCustomersByOrders(ctx context.Context, limit int) ([]*model.CustomerStats, error)
	ListTopCustomersByFrequency(ctx context.Context, limit int) ([]*model.CustomerStats, error)

	// Análisis de clientes
	ListByLevel(ctx context.Context, level string) ([]*model.CustomerStats, error)
	ListVIPCustomers(ctx context.Context) ([]*model.CustomerStats, error)
	ListInactiveCustomers(ctx context.Context, daysSince int) ([]*model.CustomerStats, error)
	ListFrequentCustomers(ctx context.Context) ([]*model.CustomerStats, error)

	// Agregaciones
	GetTotalStats(ctx context.Context) (map[string]interface{}, error)
	GetAverageOrderValue(ctx context.Context) (float64, error)
	GetTotalRevenue(ctx context.Context) (float64, error)

	// Validaciones y utilidades
	Exists(ctx context.Context, customerID int64) (bool, error)
	GetOutdatedStats(ctx context.Context) ([]*model.CustomerStats, error)
	Count(ctx context.Context) (int64, error)
}
