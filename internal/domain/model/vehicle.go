package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Vehicle representa un vehículo asociado a un cliente (para AutoParts)
type Vehicle struct {
	ID           int64            `db:"id" json:"id"`
	CustomerID   int64            `db:"customer_id" json:"customer_id" validate:"required"`
	Make         string           `db:"make" json:"make" validate:"required,min=1,max=50"`
	Model        string           `db:"model" json:"model" validate:"required,min=1,max=50"`
	Year         int              `db:"year" json:"year" validate:"required,min=1900,max=2100"`
	VIN          *string          `db:"vin" json:"vin" validate:"omitempty,max=17"`
	LicensePlate *string          `db:"license_plate" json:"license_plate" validate:"omitempty,max=20"`
	Color        *string          `db:"color" json:"color" validate:"omitempty,max=30"`
	Engine       *string          `db:"engine" json:"engine" validate:"omitempty,max=100"`
	Notes        *string          `db:"notes" json:"notes" validate:"omitempty,max=500"`
	IsActive     bool             `db:"is_active" json:"is_active"`
	Metadata     VehicleMetadata  `db:"metadata" json:"metadata"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at" json:"updated_at"`

	// Campos no persistidos (relaciones)
	Customer *Customer `db:"-" json:"customer,omitempty"`
}

// VehicleMetadata representa los metadatos del vehículo en formato JSON
type VehicleMetadata map[string]interface{}

// Implementar driver.Valuer para VehicleMetadata
func (vm VehicleMetadata) Value() (driver.Value, error) {
	if vm == nil {
		return nil, nil
	}
	return json.Marshal(vm)
}

// Implementar sql.Scanner para VehicleMetadata
func (vm *VehicleMetadata) Scan(value interface{}) error {
	if value == nil {
		*vm = make(VehicleMetadata)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al escanear VehicleMetadata: tipo inválido")
	}

	return json.Unmarshal(bytes, vm)
}

// VehicleCreate representa los datos para crear un nuevo vehículo
type VehicleCreate struct {
	CustomerID   int64
	Make         string
	Model        string
	Year         int
	VIN          *string
	LicensePlate *string
	Color        *string
	Engine       *string
	Notes        *string
	Metadata     VehicleMetadata
}

// VehicleUpdate representa los datos para actualizar un vehículo
type VehicleUpdate struct {
	ID           int64
	Make         *string
	Model        *string
	Year         *int
	VIN          *string
	LicensePlate *string
	Color        *string
	Engine       *string
	Notes        *string
	IsActive     *bool
	Metadata     VehicleMetadata
}

// VehicleFilter representa los filtros para búsqueda de vehículos
type VehicleFilter struct {
	CustomerID int64
	Search     string
	ActiveOnly bool
	Page       int
	Limit      int
}

// NewVehicle crea un nuevo vehículo desde VehicleCreate
func NewVehicle(create VehicleCreate) *Vehicle {
	now := time.Now()
	
	vehicle := &Vehicle{
		CustomerID:   create.CustomerID,
		Make:         create.Make,
		Model:        create.Model,
		Year:         create.Year,
		VIN:          create.VIN,
		LicensePlate: create.LicensePlate,
		Color:        create.Color,
		Engine:       create.Engine,
		Notes:        create.Notes,
		IsActive:     true, // Por defecto activo
		Metadata:     create.Metadata,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Asegurar que Metadata no sea nil
	if vehicle.Metadata == nil {
		vehicle.Metadata = make(VehicleMetadata)
	}

	return vehicle
}

// DisplayName devuelve una representación string del vehículo
func (v *Vehicle) DisplayName() string {
	return fmt.Sprintf("%d %s %s", v.Year, v.Make, v.Model)
}

// FullDescription devuelve una descripción completa del vehículo
func (v *Vehicle) FullDescription() string {
	description := v.DisplayName()
	
	if v.Color != nil && *v.Color != "" {
		description += fmt.Sprintf(" (%s)", *v.Color)
	}
	
	if v.LicensePlate != nil && *v.LicensePlate != "" {
		description += fmt.Sprintf(" - Placa: %s", *v.LicensePlate)
	}
	
	return description
}

// Activate activa el vehículo
func (v *Vehicle) Activate() {
	v.IsActive = true
	v.UpdatedAt = time.Now()
}

// Deactivate desactiva el vehículo
func (v *Vehicle) Deactivate() {
	v.IsActive = false
	v.UpdatedAt = time.Now()
}

// UpdateFromUpdate actualiza el vehículo con los datos de VehicleUpdate
func (v *Vehicle) UpdateFromUpdate(update VehicleUpdate) {
	now := time.Now()

	if update.Make != nil {
		v.Make = *update.Make
	}
	if update.Model != nil {
		v.Model = *update.Model
	}
	if update.Year != nil {
		v.Year = *update.Year
	}
	if update.VIN != nil {
		v.VIN = update.VIN
	}
	if update.LicensePlate != nil {
		v.LicensePlate = update.LicensePlate
	}
	if update.Color != nil {
		v.Color = update.Color
	}
	if update.Engine != nil {
		v.Engine = update.Engine
	}
	if update.Notes != nil {
		v.Notes = update.Notes
	}
	if update.IsActive != nil {
		v.IsActive = *update.IsActive
	}
	if update.Metadata != nil {
		v.Metadata = update.Metadata
	}

	v.UpdatedAt = now
}

// SetMetadata establece un valor en los metadatos
func (v *Vehicle) SetMetadata(key string, value interface{}) {
	if v.Metadata == nil {
		v.Metadata = make(VehicleMetadata)
	}
	v.Metadata[key] = value
	v.UpdatedAt = time.Now()
}

// GetMetadata obtiene un valor de los metadatos
func (v *Vehicle) GetMetadata(key string) (interface{}, bool) {
	if v.Metadata == nil {
		return nil, false
	}
	value, exists := v.Metadata[key]
	return value, exists
}

// GetMetadataString obtiene un metadato como string
func (v *Vehicle) GetMetadataString(key string) string {
	if value, exists := v.GetMetadata(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetMetadataInt obtiene un metadato como entero
func (v *Vehicle) GetMetadataInt(key string) int {
	if value, exists := v.GetMetadata(key); exists {
		switch val := value.(type) {
		case int:
			return val
		case float64:
			return int(val)
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				return i
			}
		}
	}
	return 0
}

// HasVIN verifica si el vehículo tiene VIN
func (v *Vehicle) HasVIN() bool {
	return v.VIN != nil && *v.VIN != ""
}

// HasLicensePlate verifica si el vehículo tiene placa
func (v *Vehicle) HasLicensePlate() bool {
	return v.LicensePlate != nil && *v.LicensePlate != ""
}

// IsCompatibleWith verifica compatibilidad con otro vehículo (mismo make, model, año similar)
func (v *Vehicle) IsCompatibleWith(other *Vehicle) bool {
	if other == nil {
		return false
	}

	// Mismo fabricante y modelo
	if v.Make != other.Make || v.Model != other.Model {
		return false
	}

	// Años dentro de un rango de 3 años
	yearDiff := v.Year - other.Year
	if yearDiff < 0 {
		yearDiff = -yearDiff
	}

	return yearDiff <= 3
}

// GetCompatibilityString devuelve una string de compatibilidad para búsquedas
func (v *Vehicle) GetCompatibilityString() string {
	return fmt.Sprintf("%s %s %d", v.Make, v.Model, v.Year)
}

// Validate valida los datos del vehículo
func (v *Vehicle) Validate() error {
	if v.CustomerID <= 0 {
		return &ValidationError{Field: "customer_id", Message: "ID de cliente es requerido"}
	}
	if v.Make == "" {
		return &ValidationError{Field: "make", Message: "la marca es requerida"}
	}
	if v.Model == "" {
		return &ValidationError{Field: "model", Message: "el modelo es requerido"}
	}
	if v.Year < 1900 || v.Year > 2100 {
		return &ValidationError{Field: "year", Message: "año inválido"}
	}
	if v.VIN != nil && len(*v.VIN) > 0 && len(*v.VIN) != 17 {
		return &ValidationError{Field: "vin", Message: "VIN debe tener exactamente 17 caracteres"}
	}
	return nil
}

// ValidateVIN valida que el VIN tenga el formato correcto
func (v *Vehicle) ValidateVIN() error {
	if v.VIN == nil || *v.VIN == "" {
		return nil // VIN es opcional
	}

	vin := *v.VIN
	if len(vin) != 17 {
		return &ValidationError{Field: "vin", Message: "VIN debe tener exactamente 17 caracteres"}
	}

	// Verificar que no contenga caracteres prohibidos (I, O, Q)
	prohibitedChars := "IOQ"
	for _, char := range vin {
		for _, prohibited := range prohibitedChars {
			if char == prohibited {
				return &ValidationError{Field: "vin", Message: "VIN no puede contener las letras I, O o Q"}
			}
		}
	}

	return nil
}
