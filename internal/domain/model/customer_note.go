package model

import (
	"fmt"
	"time"
)

// Constantes de tipo de nota
const (
	NoteTypeGeneral    = "general"
	NoteTypeService    = "service"
	NoteTypeComplaint  = "complaint"
	NoteTypeCompliment = "compliment"
	NoteTypeReminder   = "reminder"
	NoteTypeWarning    = "warning"
)

// CustomerNote representa una nota sobre un cliente
type CustomerNote struct {
	ID         string    `db:"id" json:"id"`
	CustomerID string    `db:"customer_id" json:"customer_id" validate:"required"`
	StaffID    string    `db:"staff_id" json:"staff_id" validate:"required"`
	StaffName  string    `db:"staff_name" json:"staff_name" validate:"required,min=1,max=200"`
	Note       string    `db:"note" json:"note" validate:"required,min=1,max=2000"`
	Type       string    `db:"type" json:"type" validate:"required,oneof=general service complaint compliment reminder warning"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`

	// Campos no persistidos (relaciones)
	Customer *Customer `db:"-" json:"customer,omitempty"`
}

// CustomerNoteCreate representa los datos para crear una nueva nota
type CustomerNoteCreate struct {
	CustomerID string
	StaffID    string
	StaffName  string
	Note       string
	Type       string
}

// CustomerNoteFilter representa los filtros para búsqueda de notas
type CustomerNoteFilter struct {
	CustomerID string
	Type       string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	Limit      int
}

// NewCustomerNote crea una nueva nota desde CustomerNoteCreate
func NewCustomerNote(create CustomerNoteCreate) *CustomerNote {
	now := time.Now()

	note := &CustomerNote{
		CustomerID: create.CustomerID,
		StaffID:    create.StaffID,
		StaffName:  create.StaffName,
		Note:       create.Note,
		Type:       create.Type,
		CreatedAt:  now,
	}

	// Si el tipo no está especificado, usar general por defecto
	if note.Type == "" {
		note.Type = NoteTypeGeneral
	}

	return note
}

// IsComplaint verifica si la nota es una queja
func (cn *CustomerNote) IsComplaint() bool {
	return cn.Type == NoteTypeComplaint
}

// IsCompliment verifica si la nota es un elogio
func (cn *CustomerNote) IsCompliment() bool {
	return cn.Type == NoteTypeCompliment
}

// IsService verifica si la nota es relacionada a servicios
func (cn *CustomerNote) IsService() bool {
	return cn.Type == NoteTypeService
}

// IsReminder verifica si la nota es un recordatorio
func (cn *CustomerNote) IsReminder() bool {
	return cn.Type == NoteTypeReminder
}

// IsWarning verifica si la nota es una advertencia
func (cn *CustomerNote) IsWarning() bool {
	return cn.Type == NoteTypeWarning
}

// GetTypeDisplayName devuelve el nombre del tipo para mostrar
func (cn *CustomerNote) GetTypeDisplayName() string {
	switch cn.Type {
	case NoteTypeGeneral:
		return "General"
	case NoteTypeService:
		return "Servicio"
	case NoteTypeComplaint:
		return "Queja"
	case NoteTypeCompliment:
		return "Elogio"
	case NoteTypeReminder:
		return "Recordatorio"
	case NoteTypeWarning:
		return "Advertencia"
	default:
		return "Desconocido"
	}
}

// GetTypeEmoji devuelve un emoji representativo del tipo
func (cn *CustomerNote) GetTypeEmoji() string {
	switch cn.Type {
	case NoteTypeGeneral:
		return "📝"
	case NoteTypeService:
		return "🔧"
	case NoteTypeComplaint:
		return "😞"
	case NoteTypeCompliment:
		return "😊"
	case NoteTypeReminder:
		return "⏰"
	case NoteTypeWarning:
		return "⚠️"
	default:
		return "📄"
	}
}

// FormattedCreatedAt devuelve la fecha de creación formateada
func (cn *CustomerNote) FormattedCreatedAt() string {
	return cn.CreatedAt.Format("02/01/2006 15:04")
}

// ShortNote devuelve una versión corta de la nota (para listados)
func (cn *CustomerNote) ShortNote(maxLength int) string {
	if len(cn.Note) <= maxLength {
		return cn.Note
	}
	return cn.Note[:maxLength-3] + "..."
}

// Summary devuelve un resumen de la nota incluyendo tipo y autor
func (cn *CustomerNote) Summary() string {
	return fmt.Sprintf("[%s] %s - %s",
		cn.GetTypeDisplayName(),
		cn.StaffName,
		cn.FormattedCreatedAt())
}

// Validate valida los datos de la nota
func (cn *CustomerNote) Validate() error {
	if cn.CustomerID == "" {
		return &ValidationError{Field: "customer_id", Message: "ID de cliente es requerido"}
	}
	if cn.StaffID == "" {
		return &ValidationError{Field: "staff_id", Message: "ID de staff es requerido"}
	}
	if cn.StaffName == "" {
		return &ValidationError{Field: "staff_name", Message: "nombre del staff es requerido"}
	}
	if cn.Note == "" {
		return &ValidationError{Field: "note", Message: "el contenido de la nota es requerido"}
	}
	if len(cn.Note) > 2000 {
		return &ValidationError{Field: "note", Message: "la nota no puede exceder 2000 caracteres"}
	}
	if !isValidNoteType(cn.Type) {
		return &ValidationError{Field: "type", Message: "tipo de nota inválido"}
	}
	return nil
}

// isValidNoteType verifica si el tipo de nota es válido
func isValidNoteType(noteType string) bool {
	validTypes := []string{
		NoteTypeGeneral,
		NoteTypeService,
		NoteTypeComplaint,
		NoteTypeCompliment,
		NoteTypeReminder,
		NoteTypeWarning,
	}

	for _, validType := range validTypes {
		if noteType == validType {
			return true
		}
	}
	return false
}

// GetValidNoteTypes devuelve todos los tipos de nota válidos
func GetValidNoteTypes() []string {
	return []string{
		NoteTypeGeneral,
		NoteTypeService,
		NoteTypeComplaint,
		NoteTypeCompliment,
		NoteTypeReminder,
		NoteTypeWarning,
	}
}

// GetNoteTypeDisplayNames devuelve un mapa de tipos a nombres de display
func GetNoteTypeDisplayNames() map[string]string {
	return map[string]string{
		NoteTypeGeneral:    "General",
		NoteTypeService:    "Servicio",
		NoteTypeComplaint:  "Queja",
		NoteTypeCompliment: "Elogio",
		NoteTypeReminder:   "Recordatorio",
		NoteTypeWarning:    "Advertencia",
	}
}
