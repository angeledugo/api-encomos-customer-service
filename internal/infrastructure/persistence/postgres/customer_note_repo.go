package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
	"github.com/encomos/api-encomos/customer-service/internal/port/repository"
)

type customerNoteRepository struct {
	db *DB
}

// NewCustomerNoteRepository creates a new customer note repository
func NewCustomerNoteRepository(db *DB) repository.CustomerNoteRepository {
	return &customerNoteRepository{
		db: db,
	}
}

// Create creates a new customer note
func (r *customerNoteRepository) Create(ctx context.Context, note *model.CustomerNote) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO customer_notes (
			customer_id, staff_id, staff_name, note, type, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id, created_at`

	err = r.db.QueryRowWithTenant(ctx, tenantID, query,
		note.CustomerID,
		note.StaffID,
		note.StaffName,
		note.Note,
		note.Type,
		note.CreatedAt,
	).Scan(&note.ID, &note.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create customer note: %w", err)
	}

	return nil
}

// GetByID retrieves a customer note by ID
func (r *customerNoteRepository) GetByID(ctx context.Context, id string) (*model.CustomerNote, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT cn.id, cn.customer_id, cn.staff_id, cn.staff_name,
			   cn.note, cn.type, cn.created_at
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.id = $1`

	note := &model.CustomerNote{}

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, id).Scan(
		&note.ID,
		&note.CustomerID,
		&note.StaffID,
		&note.StaffName,
		&note.Note,
		&note.Type,
		&note.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer note with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get customer note: %w", err)
	}

	return note, nil
}

// Delete deletes a customer note
func (r *customerNoteRepository) Delete(ctx context.Context, id string) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		DELETE FROM customer_notes
		USING customers c
		WHERE customer_notes.id = $1 AND customer_notes.customer_id = c.id`

	result, err := r.db.ExecWithTenant(ctx, tenantID, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer note with ID %s not found", id)
	}

	return nil
}

// List retrieves customer notes with filtering and pagination
func (r *customerNoteRepository) List(ctx context.Context, filter model.CustomerNoteFilter) ([]*model.CustomerNote, int, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Build WHERE clause
	var whereConditions []string
	var args []interface{}
	argCount := 0

	// Always filter by customer if provided
	if filter.CustomerID != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("cn.customer_id = $%d", argCount))
		args = append(args, filter.CustomerID)
	}

	if filter.Type != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("cn.type = $%d", argCount))
		args = append(args, filter.Type)
	}

	if filter.DateFrom != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("cn.created_at >= $%d", argCount))
		args = append(args, *filter.DateFrom)
	}

	if filter.DateTo != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("cn.created_at <= $%d", argCount))
		args = append(args, *filter.DateTo)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM customer_notes cn 
		INNER JOIN customers c ON cn.customer_id = c.id 
		%s`, whereClause)

	var total int
	err = r.db.QueryRowWithTenant(ctx, tenantID, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer notes: %w", err)
	}

	// Build pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	offset := 0
	if filter.Page > 0 {
		offset = (filter.Page - 1) * limit
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT cn.id, cn.customer_id, cn.staff_id, cn.staff_name, 
			   cn.note, cn.type, cn.created_at
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		%s
		ORDER BY cn.created_at DESC
		LIMIT %d OFFSET %d`, whereClause, limit, offset)

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customer notes: %w", err)
	}
	defer rows.Close()

	var notes []*model.CustomerNote
	for rows.Next() {
		note := &model.CustomerNote{}

		err := rows.Scan(
			&note.ID,
			&note.CustomerID,
			&note.StaffID,
			&note.StaffName,
			&note.Note,
			&note.Type,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer note: %w", err)
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over customer notes: %w", err)
	}

	return notes, total, nil
}

// ListByCustomer retrieves all notes for a customer
func (r *customerNoteRepository) ListByCustomer(ctx context.Context, customerID string) ([]*model.CustomerNote, error) {
	filter := model.CustomerNoteFilter{
		CustomerID: customerID,
		Limit:      100, // Get all notes for customer
	}
	notes, _, err := r.List(ctx, filter)
	return notes, err
}

// ListByCustomerAndType retrieves notes for a customer by type
func (r *customerNoteRepository) ListByCustomerAndType(ctx context.Context, customerID string, noteType string) ([]*model.CustomerNote, error) {
	filter := model.CustomerNoteFilter{
		CustomerID: customerID,
		Type:       noteType,
		Limit:      100,
	}
	notes, _, err := r.List(ctx, filter)
	return notes, err
}

// ListByStaff retrieves notes created by a staff member
func (r *customerNoteRepository) ListByStaff(ctx context.Context, staffID string, page, limit int) ([]*model.CustomerNote, int, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Count total records
	var total int
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM customer_notes cn 
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.staff_id = $1`, staffID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notes by staff: %w", err)
	}

	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	query := `
		SELECT cn.id, cn.customer_id, cn.staff_id, cn.staff_name, 
			   cn.note, cn.type, cn.created_at
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.staff_id = $1
		ORDER BY cn.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, staffID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notes by staff: %w", err)
	}
	defer rows.Close()

	var notes []*model.CustomerNote
	for rows.Next() {
		note := &model.CustomerNote{}

		err := rows.Scan(
			&note.ID,
			&note.CustomerID,
			&note.StaffID,
			&note.StaffName,
			&note.Note,
			&note.Type,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer note: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, total, nil
}

// ListByType retrieves notes by type
func (r *customerNoteRepository) ListByType(ctx context.Context, noteType string, page, limit int) ([]*model.CustomerNote, int, error) {
	filter := model.CustomerNoteFilter{
		Type:  noteType,
		Page:  page,
		Limit: limit,
	}
	return r.List(ctx, filter)
}

// ListRecent retrieves recent notes across all customers
func (r *customerNoteRepository) ListRecent(ctx context.Context, limit int) ([]*model.CustomerNote, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT cn.id, cn.customer_id, cn.staff_id, cn.staff_name, 
			   cn.note, cn.type, cn.created_at
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		ORDER BY cn.created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent notes: %w", err)
	}
	defer rows.Close()

	var notes []*model.CustomerNote
	for rows.Next() {
		note := &model.CustomerNote{}

		err := rows.Scan(
			&note.ID,
			&note.CustomerID,
			&note.StaffID,
			&note.StaffName,
			&note.Note,
			&note.Type,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer note: %w", err)
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// ListByDateRange retrieves notes for a customer within a date range
func (r *customerNoteRepository) ListByDateRange(ctx context.Context, customerID string, from, to *time.Time) ([]*model.CustomerNote, error) {
	filter := model.CustomerNoteFilter{
		CustomerID: customerID,
		DateFrom:   from,
		DateTo:     to,
		Limit:      1000, // Large limit for date range queries
	}
	notes, _, err := r.List(ctx, filter)
	return notes, err
}

// ListRecentByCustomer retrieves recent notes for a specific customer
func (r *customerNoteRepository) ListRecentByCustomer(ctx context.Context, customerID string, limit int) ([]*model.CustomerNote, error) {
	filter := model.CustomerNoteFilter{
		CustomerID: customerID,
		Limit:      limit,
	}
	notes, _, err := r.List(ctx, filter)
	return notes, err
}

// Count counts total customer notes
func (r *customerNoteRepository) Count(ctx context.Context) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM customer_notes cn 
		INNER JOIN customers c ON cn.customer_id = c.id`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count customer notes: %w", err)
	}

	return count, nil
}

// CountByCustomer counts notes for a specific customer
func (r *customerNoteRepository) CountByCustomer(ctx context.Context, customerID string) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*)
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.customer_id = $1`, customerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notes by customer: %w", err)
	}

	return count, nil
}

// CountByType counts notes by type
func (r *customerNoteRepository) CountByType(ctx context.Context, noteType string) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM customer_notes cn 
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.type = $1`, noteType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notes by type: %w", err)
	}

	return count, nil
}

// CountByStaff counts notes created by a staff member
func (r *customerNoteRepository) CountByStaff(ctx context.Context, staffID string) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*)
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.staff_id = $1`, staffID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notes by staff: %w", err)
	}

	return count, nil
}

// GetNoteTypesCount returns count of notes by type for a customer
func (r *customerNoteRepository) GetNoteTypesCount(ctx context.Context, customerID string) (map[string]int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT cn.type, COUNT(*) as count
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		WHERE cn.customer_id = $1
		GROUP BY cn.type
		ORDER BY count DESC`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get note types count: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var noteType string
		var count int64

		err := rows.Scan(&noteType, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note type count: %w", err)
		}

		result[noteType] = count
	}

	return result, nil
}

// GetMostActiveStaff returns staff members with most notes created
func (r *customerNoteRepository) GetMostActiveStaff(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT cn.staff_id, cn.staff_name, COUNT(*) as note_count,
			   MAX(cn.created_at) as last_note_created
		FROM customer_notes cn
		INNER JOIN customers c ON cn.customer_id = c.id
		GROUP BY cn.staff_id, cn.staff_name
		ORDER BY note_count DESC, last_note_created DESC
		LIMIT $1`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get most active staff: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var staffID string
		var staffName string
		var noteCount int64
		var lastNoteCreated time.Time

		err := rows.Scan(&staffID, &staffName, &noteCount, &lastNoteCreated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan staff activity: %w", err)
		}

		result = append(result, map[string]interface{}{
			"staff_id":          staffID,
			"staff_name":        staffName,
			"note_count":        noteCount,
			"last_note_created": lastNoteCreated,
		})
	}

	return result, nil
}
