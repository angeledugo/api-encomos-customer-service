package model

import (
	"fmt"
	"time"
)

// CustomerStats representa las estadísticas calculadas de un cliente
type CustomerStats struct {
	CustomerID        int64     `db:"customer_id" json:"customer_id"`
	TotalOrders       int32     `db:"total_orders" json:"total_orders"`
	TotalSpent        float64   `db:"total_spent" json:"total_spent"`
	AverageOrderValue float64   `db:"average_order_value" json:"average_order_value"`
	LastVisit         time.Time `db:"last_visit" json:"last_visit"`
	VisitsCount       int32     `db:"visits_count" json:"visits_count"`
	FavoriteCategory  string    `db:"favorite_category" json:"favorite_category"`
	FavoriteProducts  []string  `db:"favorite_products" json:"favorite_products"`
	CalculatedAt      time.Time `db:"calculated_at" json:"calculated_at"`

	// Campos no persistidos (relaciones)
	Customer *Customer `db:"-" json:"customer,omitempty"`
}

// CustomerStatsCreate representa los datos para crear estadísticas de cliente
type CustomerStatsCreate struct {
	CustomerID        int64
	TotalOrders       int32
	TotalSpent        float64
	AverageOrderValue float64
	LastVisit         time.Time
	VisitsCount       int32
	FavoriteCategory  string
	FavoriteProducts  []string
}

// CustomerHistoryItem representa un item del historial del cliente
type CustomerHistoryItem struct {
	ID          int64                  `json:"id"`
	Type        string                 `json:"type"` // order, appointment, note, payment
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Amount      float64                `json:"amount"`
	Status      string                 `json:"status"`
	Data        map[string]interface{} `json:"data"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CustomerHistoryFilter representa los filtros para el historial del cliente
type CustomerHistoryFilter struct {
	CustomerID int64
	Type       string // order, appointment, note, payment
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	Limit      int
}

// NewCustomerStats crea nuevas estadísticas desde CustomerStatsCreate
func NewCustomerStats(create CustomerStatsCreate) *CustomerStats {
	now := time.Now()

	stats := &CustomerStats{
		CustomerID:        create.CustomerID,
		TotalOrders:       create.TotalOrders,
		TotalSpent:        create.TotalSpent,
		AverageOrderValue: create.AverageOrderValue,
		LastVisit:         create.LastVisit,
		VisitsCount:       create.VisitsCount,
		FavoriteCategory:  create.FavoriteCategory,
		FavoriteProducts:  create.FavoriteProducts,
		CalculatedAt:      now,
	}

	// Calcular promedio si no está definido
	if stats.AverageOrderValue == 0 && stats.TotalOrders > 0 {
		stats.AverageOrderValue = stats.TotalSpent / float64(stats.TotalOrders)
	}

	return stats
}

// IsActive verifica si el cliente ha tenido actividad reciente (últimos 6 meses)
func (cs *CustomerStats) IsActive() bool {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	return cs.LastVisit.After(sixMonthsAgo)
}

// IsFrequentCustomer verifica si es un cliente frecuente (más de 10 pedidos)
func (cs *CustomerStats) IsFrequentCustomer() bool {
	return cs.TotalOrders >= 10
}

// IsHighValueCustomer verifica si es un cliente de alto valor (más de $1000 gastados)
func (cs *CustomerStats) IsHighValueCustomer() bool {
	return cs.TotalSpent >= 1000.0
}

// GetCustomerLevel devuelve el nivel del cliente basado en sus estadísticas
func (cs *CustomerStats) GetCustomerLevel() string {
	if cs.TotalSpent >= 5000 {
		return "VIP"
	} else if cs.TotalSpent >= 2000 {
		return "Premium"
	} else if cs.TotalSpent >= 500 {
		return "Gold"
	} else if cs.TotalOrders >= 5 {
		return "Silver"
	}
	return "Bronze"
}

// GetCustomerLevelEmoji devuelve un emoji para el nivel del cliente
func (cs *CustomerStats) GetCustomerLevelEmoji() string {
	switch cs.GetCustomerLevel() {
	case "VIP":
		return "💎"
	case "Premium":
		return "🥇"
	case "Gold":
		return "🥈"
	case "Silver":
		return "🥉"
	default:
		return "⭐"
	}
}

// FormattedTotalSpent devuelve el total gastado formateado como moneda
func (cs *CustomerStats) FormattedTotalSpent() string {
	return fmt.Sprintf("$%.2f", cs.TotalSpent)
}

// FormattedAverageOrderValue devuelve el promedio por pedido formateado
func (cs *CustomerStats) FormattedAverageOrderValue() string {
	return fmt.Sprintf("$%.2f", cs.AverageOrderValue)
}

// FormattedLastVisit devuelve la última visita formateada
func (cs *CustomerStats) FormattedLastVisit() string {
	return cs.LastVisit.Format("02/01/2006")
}

// DaysSinceLastVisit devuelve los días desde la última visita
func (cs *CustomerStats) DaysSinceLastVisit() int {
	return int(time.Since(cs.LastVisit).Hours() / 24)
}

// GetVisitFrequency devuelve una descripción de la frecuencia de visitas
func (cs *CustomerStats) GetVisitFrequency() string {
	daysSince := cs.DaysSinceLastVisit()

	if daysSince <= 7 {
		return "Muy frecuente"
	} else if daysSince <= 30 {
		return "Frecuente"
	} else if daysSince <= 90 {
		return "Regular"
	} else if daysSince <= 180 {
		return "Ocasional"
	}
	return "Inactivo"
}

// GetSpendingPattern devuelve un patrón de gasto
func (cs *CustomerStats) GetSpendingPattern() string {
	if cs.AverageOrderValue >= 200 {
		return "Alto valor por compra"
	} else if cs.AverageOrderValue >= 100 {
		return "Valor medio por compra"
	} else if cs.TotalOrders >= 10 {
		return "Compras frecuentes de bajo valor"
	}
	return "Comprador ocasional"
}

// HasFavoriteCategory verifica si tiene categoría favorita definida
func (cs *CustomerStats) HasFavoriteCategory() bool {
	return cs.FavoriteCategory != ""
}

// HasFavoriteProducts verifica si tiene productos favoritos
func (cs *CustomerStats) HasFavoriteProducts() bool {
	return len(cs.FavoriteProducts) > 0
}

// GetTopFavoriteProducts devuelve los primeros N productos favoritos
func (cs *CustomerStats) GetTopFavoriteProducts(n int) []string {
	if len(cs.FavoriteProducts) <= n {
		return cs.FavoriteProducts
	}
	return cs.FavoriteProducts[:n]
}

// IsStatsOutdated verifica si las estadísticas están desactualizadas (más de 24 horas)
func (cs *CustomerStats) IsStatsOutdated() bool {
	return time.Since(cs.CalculatedAt) > 24*time.Hour
}

// UpdateCalculatedAt actualiza la fecha de cálculo
func (cs *CustomerStats) UpdateCalculatedAt() {
	cs.CalculatedAt = time.Now()
}

// RecalculateAverageOrderValue recalcula el valor promedio por pedido
func (cs *CustomerStats) RecalculateAverageOrderValue() {
	if cs.TotalOrders > 0 {
		cs.AverageOrderValue = cs.TotalSpent / float64(cs.TotalOrders)
	} else {
		cs.AverageOrderValue = 0
	}
}

// AddOrder actualiza las estadísticas con un nuevo pedido
func (cs *CustomerStats) AddOrder(amount float64, visitDate time.Time) {
	cs.TotalOrders++
	cs.TotalSpent += amount
	cs.VisitsCount++

	if visitDate.After(cs.LastVisit) {
		cs.LastVisit = visitDate
	}

	cs.RecalculateAverageOrderValue()
	cs.UpdateCalculatedAt()
}

// Validate valida las estadísticas del cliente
func (cs *CustomerStats) Validate() error {
	if cs.CustomerID <= 0 {
		return &ValidationError{Field: "customer_id", Message: "ID de cliente es requerido"}
	}
	if cs.TotalOrders < 0 {
		return &ValidationError{Field: "total_orders", Message: "total de pedidos no puede ser negativo"}
	}
	if cs.TotalSpent < 0 {
		return &ValidationError{Field: "total_spent", Message: "total gastado no puede ser negativo"}
	}
	if cs.VisitsCount < 0 {
		return &ValidationError{Field: "visits_count", Message: "número de visitas no puede ser negativo"}
	}
	return nil
}

// GetSummary devuelve un resumen de las estadísticas
func (cs *CustomerStats) GetSummary() string {
	return fmt.Sprintf("%s - %d pedidos, %s gastado, última visita: %s",
		cs.GetCustomerLevel(),
		cs.TotalOrders,
		cs.FormattedTotalSpent(),
		cs.FormattedLastVisit())
}
