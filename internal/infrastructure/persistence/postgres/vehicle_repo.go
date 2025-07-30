package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourorg/api-encomos/customer-service/internal/domain/model"
	"github.com/yourorg/api-encomos/customer-service/internal/port/repository"
)

type vehicleRepository struct {
	db *DB
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(db *DB) repository.VehicleRepository {
	return &vehicleRepository{
		db: db,
	}
}

// Create creates a new vehicle
func (r *vehicleRepository) Create(ctx context.Context, vehicle *model.Vehicle) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO vehicles (
			customer_id, make, model, year, vin, license_plate,
			color, engine, notes, is_active, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at`

	err = r.db.QueryRowWithTenant(ctx, tenantID, query,
		vehicle.CustomerID,
		vehicle.Make,
		vehicle.Model,
		vehicle.Year,
		NullString(vehicle.VIN),
		NullString(vehicle.LicensePlate),
		NullString(vehicle.Color),
		NullString(vehicle.Engine),
		NullString(vehicle.Notes),
		vehicle.IsActive,
		vehicle.Metadata,
		vehicle.CreatedAt,
		vehicle.UpdatedAt,
	).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	return nil
}

// GetByID retrieves a vehicle by ID
func (r *vehicleRepository) GetByID(ctx context.Context, id int64) (*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.id = $1`

	vehicle := &model.Vehicle{}
	var vin, licensePlate, color, engine, notes sql.NullString

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, id).Scan(
		&vehicle.ID,
		&vehicle.CustomerID,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.Year,
		&vin,
		&licensePlate,
		&color,
		&engine,
		&notes,
		&vehicle.IsActive,
		&vehicle.Metadata,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	// Convert nullable fields
	vehicle.VIN = StringFromNull(vin)
	vehicle.LicensePlate = StringFromNull(licensePlate)
	vehicle.Color = StringFromNull(color)
	vehicle.Engine = StringFromNull(engine)
	vehicle.Notes = StringFromNull(notes)

	return vehicle, nil
}

// Update updates a vehicle
func (r *vehicleRepository) Update(ctx context.Context, vehicle *model.Vehicle) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE vehicles SET
			make = $2, model = $3, year = $4, vin = $5,
			license_plate = $6, color = $7, engine = $8, notes = $9,
			is_active = $10, metadata = $11, updated_at = $12
		FROM customers c
		WHERE vehicles.id = $1 AND vehicles.customer_id = c.id`

	result, err := r.db.ExecWithTenant(ctx, tenantID, query,
		vehicle.ID,
		vehicle.Make,
		vehicle.Model,
		vehicle.Year,
		NullString(vehicle.VIN),
		NullString(vehicle.LicensePlate),
		NullString(vehicle.Color),
		NullString(vehicle.Engine),
		NullString(vehicle.Notes),
		vehicle.IsActive,
		vehicle.Metadata,
		vehicle.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle with ID %d not found", vehicle.ID)
	}

	return nil
}

// Delete deletes a vehicle
func (r *vehicleRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		DELETE FROM vehicles 
		USING customers c
		WHERE vehicles.id = $1 AND vehicles.customer_id = c.id`

	result, err := r.db.ExecWithTenant(ctx, tenantID, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle with ID %d not found", id)
	}

	return nil
}

// List retrieves vehicles with filtering and pagination
func (r *vehicleRepository) List(ctx context.Context, filter model.VehicleFilter) ([]*model.Vehicle, int, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Build WHERE clause
	var whereConditions []string
	var args []interface{}
	argCount := 0

	// Always filter by customer if provided
	if filter.CustomerID > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("v.customer_id = $%d", argCount))
		args = append(args, filter.CustomerID)
	}

	if filter.Search != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(v.make ILIKE $%d OR v.model ILIKE $%d OR v.vin ILIKE $%d OR v.license_plate ILIKE $%d)",
			argCount, argCount, argCount, argCount))
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.ActiveOnly {
		whereConditions = append(whereConditions, "v.is_active = true")
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id 
		%s`, whereClause)

	var total int
	err = r.db.QueryRowWithTenant(ctx, tenantID, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
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
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		%s
		ORDER BY v.year DESC, v.make, v.model
		LIMIT %d OFFSET %d`, whereClause, limit, offset)

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*model.Vehicle
	for rows.Next() {
		vehicle := &model.Vehicle{}
		var vin, licensePlate, color, engine, notes sql.NullString

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.CustomerID,
			&vehicle.Make,
			&vehicle.Model,
			&vehicle.Year,
			&vin,
			&licensePlate,
			&color,
			&engine,
			&notes,
			&vehicle.IsActive,
			&vehicle.Metadata,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		// Convert nullable fields
		vehicle.VIN = StringFromNull(vin)
		vehicle.LicensePlate = StringFromNull(licensePlate)
		vehicle.Color = StringFromNull(color)
		vehicle.Engine = StringFromNull(engine)
		vehicle.Notes = StringFromNull(notes)

		vehicles = append(vehicles, vehicle)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over vehicles: %w", err)
	}

	return vehicles, total, nil
}

// ListByCustomer retrieves all vehicles for a customer
func (r *vehicleRepository) ListByCustomer(ctx context.Context, customerID int64) ([]*model.Vehicle, error) {
	filter := model.VehicleFilter{
		CustomerID: customerID,
		Limit:      100, // Get all vehicles for customer
	}
	vehicles, _, err := r.List(ctx, filter)
	return vehicles, err
}

// GetByVIN retrieves a vehicle by VIN
func (r *vehicleRepository) GetByVIN(ctx context.Context, vin string) (*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.vin = $1`

	vehicle := &model.Vehicle{}
	var vinNull, licensePlate, color, engine, notes sql.NullString

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, vin).Scan(
		&vehicle.ID,
		&vehicle.CustomerID,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.Year,
		&vinNull,
		&licensePlate,
		&color,
		&engine,
		&notes,
		&vehicle.IsActive,
		&vehicle.Metadata,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle with VIN %s not found", vin)
		}
		return nil, fmt.Errorf("failed to get vehicle by VIN: %w", err)
	}

	// Convert nullable fields
	vehicle.VIN = StringFromNull(vinNull)
	vehicle.LicensePlate = StringFromNull(licensePlate)
	vehicle.Color = StringFromNull(color)
	vehicle.Engine = StringFromNull(engine)
	vehicle.Notes = StringFromNull(notes)

	return vehicle, nil
}

// GetByLicensePlate retrieves a vehicle by license plate
func (r *vehicleRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.license_plate = $1`

	vehicle := &model.Vehicle{}
	var vin, licensePlateNull, color, engine, notes sql.NullString

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, licensePlate).Scan(
		&vehicle.ID,
		&vehicle.CustomerID,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.Year,
		&vin,
		&licensePlateNull,
		&color,
		&engine,
		&notes,
		&vehicle.IsActive,
		&vehicle.Metadata,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle with license plate %s not found", licensePlate)
		}
		return nil, fmt.Errorf("failed to get vehicle by license plate: %w", err)
	}

	// Convert nullable fields
	vehicle.VIN = StringFromNull(vin)
	vehicle.LicensePlate = StringFromNull(licensePlateNull)
	vehicle.Color = StringFromNull(color)
	vehicle.Engine = StringFromNull(engine)
	vehicle.Notes = StringFromNull(notes)

	return vehicle, nil
}

// SearchByMakeModel searches vehicles by make and model
func (r *vehicleRepository) SearchByMakeModel(ctx context.Context, make, model string, year *int) ([]*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var whereConditions []string
	var args []interface{}
	argCount := 0

	if make != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("v.make ILIKE $%d", argCount))
		args = append(args, "%"+make+"%")
	}

	if model != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("v.model ILIKE $%d", argCount))
		args = append(args, "%"+model+"%")
	}

	if year != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("v.year = $%d", argCount))
		args = append(args, *year)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		%s
		ORDER BY v.year DESC, v.make, v.model
		LIMIT 50`, whereClause)

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*model.Vehicle
	for rows.Next() {
		vehicle := &model.Vehicle{}
		var vin, licensePlate, color, engine, notes sql.NullString

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.CustomerID,
			&vehicle.Make,
			&vehicle.Model,
			&vehicle.Year,
			&vin,
			&licensePlate,
			&color,
			&engine,
			&notes,
			&vehicle.IsActive,
			&vehicle.Metadata,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		// Convert nullable fields
		vehicle.VIN = StringFromNull(vin)
		vehicle.LicensePlate = StringFromNull(licensePlate)
		vehicle.Color = StringFromNull(color)
		vehicle.Engine = StringFromNull(engine)
		vehicle.Notes = StringFromNull(notes)

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// FindCompatibleVehicles finds vehicles compatible within a year range
func (r *vehicleRepository) FindCompatibleVehicles(ctx context.Context, make, model string, yearFrom, yearTo int) ([]*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.make ILIKE $1 AND v.model ILIKE $2 
		  AND v.year BETWEEN $3 AND $4
		  AND v.is_active = true
		ORDER BY v.year DESC
		LIMIT 100`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, "%"+make+"%", "%"+model+"%", yearFrom, yearTo)
	if err != nil {
		return nil, fmt.Errorf("failed to find compatible vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*model.Vehicle
	for rows.Next() {
		vehicle := &model.Vehicle{}
		var vin, licensePlate, color, engine, notes sql.NullString

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.CustomerID,
			&vehicle.Make,
			&vehicle.Model,
			&vehicle.Year,
			&vin,
			&licensePlate,
			&color,
			&engine,
			&notes,
			&vehicle.IsActive,
			&vehicle.Metadata,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		// Convert nullable fields
		vehicle.VIN = StringFromNull(vin)
		vehicle.LicensePlate = StringFromNull(licensePlate)
		vehicle.Color = StringFromNull(color)
		vehicle.Engine = StringFromNull(engine)
		vehicle.Notes = StringFromNull(notes)

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// ListByMakeModelYear lists vehicles by exact make, model and year
func (r *vehicleRepository) ListByMakeModelYear(ctx context.Context, make, model string, year int) ([]*model.Vehicle, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT v.id, v.customer_id, v.make, v.model, v.year, v.vin,
			   v.license_plate, v.color, v.engine, v.notes, v.is_active,
			   v.metadata, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.make = $1 AND v.model = $2 AND v.year = $3
		ORDER BY v.created_at DESC`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, make, model, year)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles by make/model/year: %w", err)
	}
	defer rows.Close()

	var vehicles []*model.Vehicle
	for rows.Next() {
		vehicle := &model.Vehicle{}
		var vin, licensePlate, color, engine, notes sql.NullString

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.CustomerID,
			&vehicle.Make,
			&vehicle.Model,
			&vehicle.Year,
			&vin,
			&licensePlate,
			&color,
			&engine,
			&notes,
			&vehicle.IsActive,
			&vehicle.Metadata,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		// Convert nullable fields
		vehicle.VIN = StringFromNull(vin)
		vehicle.LicensePlate = StringFromNull(licensePlate)
		vehicle.Color = StringFromNull(color)
		vehicle.Engine = StringFromNull(engine)
		vehicle.Notes = StringFromNull(notes)

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// Count counts total vehicles
func (r *vehicleRepository) Count(ctx context.Context) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count vehicles: %w", err)
	}

	return count, nil
}

// CountByCustomer counts vehicles for a specific customer
func (r *vehicleRepository) CountByCustomer(ctx context.Context, customerID int64) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.customer_id = $1`, customerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count vehicles by customer: %w", err)
	}

	return count, nil
}

// CountActive counts active vehicles
func (r *vehicleRepository) CountActive(ctx context.Context) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID, `
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.is_active = true`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active vehicles: %w", err)
	}

	return count, nil
}

// ExistsByVIN checks if a vehicle exists by VIN
func (r *vehicleRepository) ExistsByVIN(ctx context.Context, vin string, excludeID *int64) (bool, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := `
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.vin = $1`
	args := []interface{}{vin}

	if excludeID != nil {
		query += " AND v.id != $2"
		args = append(args, *excludeID)
	}

	var count int
	err = r.db.QueryRowWithTenant(ctx, tenantID, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check VIN existence: %w", err)
	}

	return count > 0, nil
}

// ExistsByLicensePlate checks if a vehicle exists by license plate
func (r *vehicleRepository) ExistsByLicensePlate(ctx context.Context, licensePlate string, excludeID *int64) (bool, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := `
		SELECT COUNT(*) 
		FROM vehicles v 
		INNER JOIN customers c ON v.customer_id = c.id
		WHERE v.license_plate = $1`
	args := []interface{}{licensePlate}

	if excludeID != nil {
		query += " AND v.id != $2"
		args = append(args, *excludeID)
	}

	var count int
	err = r.db.QueryRowWithTenant(ctx, tenantID, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check license plate existence: %w", err)
	}

	return count > 0, nil
}

// CreateBatch creates multiple vehicles in a transaction
func (r *vehicleRepository) CreateBatch(ctx context.Context, vehicles []*model.Vehicle) error {
	if len(vehicles) == 0 {
		return nil
	}

	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	return r.db.TransactionWithTenant(ctx, tenantID, func(tx *sql.Tx) error {
		query := `
			INSERT INTO vehicles (
				customer_id, make, model, year, vin, license_plate,
				color, engine, notes, is_active, metadata, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			) RETURNING id, created_at, updated_at`

		for _, vehicle := range vehicles {
			err := tx.QueryRowContext(ctx, query,
				vehicle.CustomerID,
				vehicle.Make,
				vehicle.Model,
				vehicle.Year,
				NullString(vehicle.VIN),
				NullString(vehicle.LicensePlate),
				NullString(vehicle.Color),
				NullString(vehicle.Engine),
				NullString(vehicle.Notes),
				vehicle.IsActive,
				vehicle.Metadata,
				vehicle.CreatedAt,
				vehicle.UpdatedAt,
			).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)

			if err != nil {
				return fmt.Errorf("failed to create vehicle in batch: %w", err)
			}
		}

		return nil
	})
}

// ListActiveByCustomer retrieves all active vehicles for a customer
func (r *vehicleRepository) ListActiveByCustomer(ctx context.Context, customerID int64) ([]*model.Vehicle, error) {
	filter := model.VehicleFilter{
		CustomerID: customerID,
		ActiveOnly: true,
		Limit:      100,
	}
	vehicles, _, err := r.List(ctx, filter)
	return vehicles, err
}
