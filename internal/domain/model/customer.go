package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Constantes de tipo de cliente
const (
	CustomerTypeIndividual = "individual"
	CustomerTypeBusiness   = "business"
)

// Customer representa un cliente en el sistema
type Customer struct {
	ID           int64               `db:"id" json:"id"`
	TenantID     int64               `db:"tenant_id" json:"tenant_id"`
	FirstName    string              `db:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName     string              `db:"last_name" json:"last_name" validate:"required,min=1,max=100"`
	Email        *string             `db:"email" json:"email" validate:"omitempty,email,max=255"`
	Phone        *string             `db:"phone" json:"phone" validate:"omitempty,max=20"`
	CustomerType string              `db:"customer_type" json:"customer_type" validate:"required,oneof=individual business"`
	CompanyName  *string             `db:"company_name" json:"company_name" validate:"omitempty,max=255"`
	TaxID        *string             `db:"tax_id" json:"tax_id" validate:"omitempty,max=50"`
	Address      *string             `db:"address" json:"address" validate:"omitempty,max=500"`
	Birthday     *time.Time          `db:"birthday" json:"birthday"`
	Notes        *string             `db:"notes" json:"notes" validate:"omitempty,max=1000"`
	Preferences  CustomerPreferences `db:"preferences" json:"preferences"`
	IsActive     bool                `db:"is_active" json:"is_active"`
	CreatedAt    time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updated_at"`

	// Campos no persistidos (relaciones)
	Vehicles      []*Vehicle      `db:"-" json:"vehicles,omitempty"`
	CustomerNotes []*CustomerNote `db:"-" json:"customer_notes,omitempty"`
	Stats         *CustomerStats  `db:"-" json:"stats,omitempty"`
}

// CustomerPreferences representa las preferencias del cliente en formato JSON
type CustomerPreferences map[string]interface{}

// Implementar driver.Valuer para CustomerPreferences
func (cp CustomerPreferences) Value() (driver.Value, error) {
	if cp == nil {
		return nil, nil
	}
	return json.Marshal(cp)
}

// Implementar sql.Scanner para CustomerPreferences
func (cp *CustomerPreferences) Scan(value interface{}) error {
	if value == nil {
		*cp = make(CustomerPreferences)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error al escanear CustomerPreferences: tipo inválido")
	}

	return json.Unmarshal(bytes, cp)
}

// CustomerCreate representa los datos para crear un nuevo cliente
type CustomerCreate struct {
	TenantID     int64
	FirstName    string
	LastName     string
	Email        *string
	Phone        *string
	CustomerType string
	CompanyName  *string
	TaxID        *string
	Address      *string
	Birthday     *time.Time
	Notes        *string
	Preferences  CustomerPreferences
}

// CustomerUpdate representa los datos para actualizar un cliente
type CustomerUpdate struct {
	ID           int64
	FirstName    *string
	LastName     *string
	Email        *string
	Phone        *string
	CustomerType *string
	CompanyName  *string
	TaxID        *string
	Address      *string
	Birthday     *time.Time
	Notes        *string
	Preferences  CustomerPreferences
	IsActive     *bool
}

// CustomerFilter representa los filtros para búsqueda de clientes
type CustomerFilter struct {
	Search       string
	CustomerType string
	ActiveOnly   bool
	Page         int
	Limit        int
	SortBy       string // name, created_at, last_visit, total_spent
	SortOrder    string // asc, desc
}

// CustomerSearchFilter representa los filtros para búsqueda avanzada
type CustomerSearchFilter struct {
	Query        string
	SearchFields []string // name, email, phone, tax_id
	Limit        int
}

// NewCustomer crea un nuevo cliente desde CustomerCreate
func NewCustomer(create CustomerCreate) *Customer {
	now := time.Now()

	customer := &Customer{
		TenantID:     create.TenantID,
		FirstName:    create.FirstName,
		LastName:     create.LastName,
		Email:        create.Email,
		Phone:        create.Phone,
		CustomerType: create.CustomerType,
		CompanyName:  create.CompanyName,
		TaxID:        create.TaxID,
		Address:      create.Address,
		Birthday:     create.Birthday,
		Notes:        create.Notes,
		Preferences:  create.Preferences,
		IsActive:     true, // Por defecto activo
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Asegurar que Preferences no sea nil
	if customer.Preferences == nil {
		customer.Preferences = make(CustomerPreferences)
	}

	return customer
}

// FullName devuelve el nombre completo del cliente
func (c *Customer) FullName() string {
	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

// DisplayName devuelve el nombre para mostrar (nombre completo o empresa)
func (c *Customer) DisplayName() string {
	if c.CustomerType == CustomerTypeBusiness && c.CompanyName != nil && *c.CompanyName != "" {
		return *c.CompanyName
	}
	return c.FullName()
}

// IsIndividual verifica si es un cliente individual
func (c *Customer) IsIndividual() bool {
	return c.CustomerType == CustomerTypeIndividual
}

// IsBusiness verifica si es un cliente empresa
func (c *Customer) IsBusiness() bool {
	return c.CustomerType == CustomerTypeBusiness
}

// Activate activa el cliente
func (c *Customer) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate desactiva el cliente
func (c *Customer) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// UpdateFromUpdate actualiza el cliente con los datos de CustomerUpdate
func (c *Customer) UpdateFromUpdate(update CustomerUpdate) {
	now := time.Now()

	if update.FirstName != nil {
		c.FirstName = *update.FirstName
	}
	if update.LastName != nil {
		c.LastName = *update.LastName
	}
	if update.Email != nil {
		c.Email = update.Email
	}
	if update.Phone != nil {
		c.Phone = update.Phone
	}
	if update.CustomerType != nil {
		c.CustomerType = *update.CustomerType
	}
	if update.CompanyName != nil {
		c.CompanyName = update.CompanyName
	}
	if update.TaxID != nil {
		c.TaxID = update.TaxID
	}
	if update.Address != nil {
		c.Address = update.Address
	}
	if update.Birthday != nil {
		c.Birthday = update.Birthday
	}
	if update.Notes != nil {
		c.Notes = update.Notes
	}
	if update.Preferences != nil {
		c.Preferences = update.Preferences
	}
	if update.IsActive != nil {
		c.IsActive = *update.IsActive
	}

	c.UpdatedAt = now
}

// SetPreference establece una preferencia específica
func (c *Customer) SetPreference(key string, value interface{}) {
	if c.Preferences == nil {
		c.Preferences = make(CustomerPreferences)
	}
	c.Preferences[key] = value
	c.UpdatedAt = time.Now()
}

// GetPreference obtiene una preferencia específica
func (c *Customer) GetPreference(key string) (interface{}, bool) {
	if c.Preferences == nil {
		return nil, false
	}
	value, exists := c.Preferences[key]
	return value, exists
}

// GetPreferenceString obtiene una preferencia como string
func (c *Customer) GetPreferenceString(key string) string {
	if value, exists := c.GetPreference(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetPreferenceBool obtiene una preferencia como boolean
func (c *Customer) GetPreferenceBool(key string) bool {
	if value, exists := c.GetPreference(key); exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

// HasEmail verifica si el cliente tiene email
func (c *Customer) HasEmail() bool {
	return c.Email != nil && *c.Email != ""
}

// HasPhone verifica si el cliente tiene teléfono
func (c *Customer) HasPhone() bool {
	return c.Phone != nil && *c.Phone != ""
}

// HasBirthday verifica si el cliente tiene fecha de cumpleaños
func (c *Customer) HasBirthday() bool {
	return c.Birthday != nil
}

// Age calcula la edad del cliente
func (c *Customer) Age() *int {
	if c.Birthday == nil {
		return nil
	}
	age := int(time.Since(*c.Birthday).Hours() / 24 / 365)
	return &age
}

// Validate valida los datos del cliente
func (c *Customer) Validate() error {
	if c.FirstName == "" {
		return &ValidationError{Field: "first_name", Message: "el nombre es requerido"}
	}
	if c.LastName == "" {
		return &ValidationError{Field: "last_name", Message: "el apellido es requerido"}
	}
	if c.CustomerType != CustomerTypeIndividual && c.CustomerType != CustomerTypeBusiness {
		return &ValidationError{Field: "customer_type", Message: "tipo de cliente inválido"}
	}
	if c.CustomerType == CustomerTypeBusiness && (c.CompanyName == nil || *c.CompanyName == "") {
		return &ValidationError{Field: "company_name", Message: "el nombre de la empresa es requerido para clientes empresariales"}
	}
	if c.Email != nil && *c.Email != "" {
		// Validación básica de email
		if !isValidEmail(*c.Email) {
			return &ValidationError{Field: "email", Message: "formato de email inválido"}
		}
	}
	return nil
}

// isValidEmail realiza una validación básica de email
func isValidEmail(email string) bool {
	// Implementación simple, en producción usar una validación más robusta
	return len(email) > 3 &&
		len(email) <= 255 &&
		containsChar(email, '@') &&
		containsChar(email, '.')
}

// containsChar verifica si un string contiene un carácter específico
func containsChar(s string, char rune) bool {
	for _, c := range s {
		if c == char {
			return true
		}
	}
	return false
}

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
