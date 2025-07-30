package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/encomos/api-encomos/customer-service/internal/domain/model"
	"github.com/encomos/api-encomos/customer-service/internal/port/repository"
)

type customerRepository struct {
	db *DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *DB) repository.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// Create creates a new customer
func (r *customerRepository) Create(ctx context.Context, customer *model.Customer) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO customers (
			tenant_id, first_name, last_name, email, phone, 
			customer_type, company_name, tax_id, address, birthday, 
			notes, preferences, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id, created_at, updated_at`

	err = r.db.QueryRowWithTenant(ctx, tenantID, query,
		tenantID,
		customer.FirstName,
		customer.LastName,
		NullString(customer.Email),
		NullString(customer.Phone),
		customer.CustomerType,
		NullString(customer.CompanyName),
		NullString(customer.TaxID),
		NullString(customer.Address),
		NullTime(customer.Birthday),
		NullString(customer.Notes),
		customer.Preferences,
		customer.IsActive,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	customer.TenantID = tenantID
	return nil
}

// GetByID retrieves a customer by ID
func (r *customerRepository) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		WHERE id = $1`

	customer := &model.Customer{}
	var email, phone, companyName, taxID, address, notes sql.NullString
	var birthday sql.NullTime

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, id).Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.FirstName,
		&customer.LastName,
		&email,
		&phone,
		&customer.CustomerType,
		&companyName,
		&taxID,
		&address,
		&birthday,
		&notes,
		&customer.Preferences,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Convert nullable fields
	customer.Email = StringFromNull(email)
	customer.Phone = StringFromNull(phone)
	customer.CompanyName = StringFromNull(companyName)
	customer.TaxID = StringFromNull(taxID)
	customer.Address = StringFromNull(address)
	customer.Notes = StringFromNull(notes)
	customer.Birthday = TimeFromNull(birthday)

	return customer, nil
}

// Update updates a customer
func (r *customerRepository) Update(ctx context.Context, customer *model.Customer) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE customers SET
			first_name = $2, last_name = $3, email = $4, phone = $5,
			customer_type = $6, company_name = $7, tax_id = $8, address = $9,
			birthday = $10, notes = $11, preferences = $12, is_active = $13,
			updated_at = $14
		WHERE id = $1`

	result, err := r.db.ExecWithTenant(ctx, tenantID, query,
		customer.ID,
		customer.FirstName,
		customer.LastName,
		NullString(customer.Email),
		NullString(customer.Phone),
		customer.CustomerType,
		NullString(customer.CompanyName),
		NullString(customer.TaxID),
		NullString(customer.Address),
		NullTime(customer.Birthday),
		NullString(customer.Notes),
		customer.Preferences,
		customer.IsActive,
		customer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer with ID %d not found", customer.ID)
	}

	return nil
}

// Delete deletes a customer
func (r *customerRepository) Delete(ctx context.Context, id int64) error {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM customers WHERE id = $1`

	result, err := r.db.ExecWithTenant(ctx, tenantID, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer with ID %d not found", id)
	}

	return nil
}

// List retrieves customers with filtering and pagination
func (r *customerRepository) List(ctx context.Context, filter model.CustomerFilter) ([]*model.Customer, int, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Build WHERE clause
	var whereConditions []string
	var args []interface{}
	argCount := 0

	if filter.Search != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR company_name ILIKE $%d)",
			argCount, argCount, argCount, argCount))
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.CustomerType != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("customer_type = $%d", argCount))
		args = append(args, filter.CustomerType)
	}

	if filter.ActiveOnly {
		whereConditions = append(whereConditions, "is_active = true")
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Build ORDER BY clause
	orderBy := "ORDER BY created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}

		switch filter.SortBy {
		case "name":
			orderBy = fmt.Sprintf("ORDER BY first_name %s, last_name %s", direction, direction)
		case "created_at":
			orderBy = fmt.Sprintf("ORDER BY created_at %s", direction)
		case "company_name":
			orderBy = fmt.Sprintf("ORDER BY company_name %s", direction)
		}
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", whereClause)
	var total int
	err = r.db.QueryRowWithTenant(ctx, tenantID, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
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
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		%s %s
		LIMIT %d OFFSET %d`, whereClause, orderBy, limit, offset)

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		customer := &model.Customer{}
		var email, phone, companyName, taxID, address, notes sql.NullString
		var birthday sql.NullTime

		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.FirstName,
			&customer.LastName,
			&email,
			&phone,
			&customer.CustomerType,
			&companyName,
			&taxID,
			&address,
			&birthday,
			&notes,
			&customer.Preferences,
			&customer.IsActive,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer: %w", err)
		}

		// Convert nullable fields
		customer.Email = StringFromNull(email)
		customer.Phone = StringFromNull(phone)
		customer.CompanyName = StringFromNull(companyName)
		customer.TaxID = StringFromNull(taxID)
		customer.Address = StringFromNull(address)
		customer.Notes = StringFromNull(notes)
		customer.Birthday = TimeFromNull(birthday)

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over customers: %w", err)
	}

	return customers, total, nil
}

// Search performs advanced search on customers
func (r *customerRepository) Search(ctx context.Context, filter model.CustomerSearchFilter) ([]*model.Customer, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if filter.Query == "" {
		return []*model.Customer{}, nil
	}

	// Build search conditions based on search fields
	var searchConditions []string
	searchFields := filter.SearchFields
	if len(searchFields) == 0 {
		// Default search fields
		searchFields = []string{"name", "email", "phone", "tax_id"}
	}

	for _, field := range searchFields {
		switch field {
		case "name":
			searchConditions = append(searchConditions,
				"(first_name ILIKE $1 OR last_name ILIKE $1 OR (first_name || ' ' || last_name) ILIKE $1)")
		case "email":
			searchConditions = append(searchConditions, "email ILIKE $1")
		case "phone":
			searchConditions = append(searchConditions, "phone ILIKE $1")
		case "tax_id":
			searchConditions = append(searchConditions, "tax_id ILIKE $1")
		case "company_name":
			searchConditions = append(searchConditions, "company_name ILIKE $1")
		}
	}

	if len(searchConditions) == 0 {
		return []*model.Customer{}, nil
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20 // Default limit for search
	}

	query := fmt.Sprintf(`
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		WHERE (%s) AND is_active = true
		ORDER BY 
			CASE 
				WHEN first_name ILIKE $1 OR last_name ILIKE $1 THEN 1
				WHEN email = $2 THEN 2
				WHEN phone = $2 THEN 3
				ELSE 4
			END,
			first_name, last_name
		LIMIT %d`, strings.Join(searchConditions, " OR "), limit)

	searchTerm := "%" + filter.Query + "%"
	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, searchTerm, filter.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		customer := &model.Customer{}
		var email, phone, companyName, taxID, address, notes sql.NullString
		var birthday sql.NullTime

		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.FirstName,
			&customer.LastName,
			&email,
			&phone,
			&customer.CustomerType,
			&companyName,
			&taxID,
			&address,
			&birthday,
			&notes,
			&customer.Preferences,
			&customer.IsActive,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		// Convert nullable fields
		customer.Email = StringFromNull(email)
		customer.Phone = StringFromNull(phone)
		customer.CompanyName = StringFromNull(companyName)
		customer.TaxID = StringFromNull(taxID)
		customer.Address = StringFromNull(address)
		customer.Notes = StringFromNull(notes)
		customer.Birthday = TimeFromNull(birthday)

		customers = append(customers, customer)
	}

	return customers, nil
}

// GetByEmail retrieves a customer by email
func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		WHERE email = $1`

	customer := &model.Customer{}
	var emailNull, phone, companyName, taxID, address, notes sql.NullString
	var birthday sql.NullTime

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, email).Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.FirstName,
		&customer.LastName,
		&emailNull,
		&phone,
		&customer.CustomerType,
		&companyName,
		&taxID,
		&address,
		&birthday,
		&notes,
		&customer.Preferences,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get customer by email: %w", err)
	}

	// Convert nullable fields
	customer.Email = StringFromNull(emailNull)
	customer.Phone = StringFromNull(phone)
	customer.CompanyName = StringFromNull(companyName)
	customer.TaxID = StringFromNull(taxID)
	customer.Address = StringFromNull(address)
	customer.Notes = StringFromNull(notes)
	customer.Birthday = TimeFromNull(birthday)

	return customer, nil
}

// GetByTaxID retrieves a customer by tax ID
func (r *customerRepository) GetByTaxID(ctx context.Context, taxID string) (*model.Customer, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		WHERE tax_id = $1`

	customer := &model.Customer{}
	var email, phone, companyName, taxIDNull, address, notes sql.NullString
	var birthday sql.NullTime

	err = r.db.QueryRowWithTenant(ctx, tenantID, query, taxID).Scan(
		&customer.ID,
		&customer.TenantID,
		&customer.FirstName,
		&customer.LastName,
		&email,
		&phone,
		&customer.CustomerType,
		&companyName,
		&taxIDNull,
		&address,
		&birthday,
		&notes,
		&customer.Preferences,
		&customer.IsActive,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with tax ID %s not found", taxID)
		}
		return nil, fmt.Errorf("failed to get customer by tax ID: %w", err)
	}

	// Convert nullable fields
	customer.Email = StringFromNull(email)
	customer.Phone = StringFromNull(phone)
	customer.CompanyName = StringFromNull(companyName)
	customer.TaxID = StringFromNull(taxIDNull)
	customer.Address = StringFromNull(address)
	customer.Notes = StringFromNull(notes)
	customer.Birthday = TimeFromNull(birthday)

	return customer, nil
}

// ListByType retrieves customers by type with pagination
func (r *customerRepository) ListByType(ctx context.Context, customerType string, page, limit int) ([]*model.Customer, int, error) {
	filter := model.CustomerFilter{
		CustomerType: customerType,
		Page:         page,
		Limit:        limit,
	}
	return r.List(ctx, filter)
}

// ListActive retrieves active customers with pagination
func (r *customerRepository) ListActive(ctx context.Context, page, limit int) ([]*model.Customer, int, error) {
	filter := model.CustomerFilter{
		ActiveOnly: true,
		Page:       page,
		Limit:      limit,
	}
	return r.List(ctx, filter)
}

// ListInactive retrieves inactive customers with pagination
func (r *customerRepository) ListInactive(ctx context.Context, page, limit int) ([]*model.Customer, int, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Count total inactive customers
	var total int
	err = r.db.QueryRowWithTenant(ctx, tenantID,
		"SELECT COUNT(*) FROM customers WHERE is_active = false").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inactive customers: %w", err)
	}

	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	query := `
		SELECT id, tenant_id, first_name, last_name, email, phone,
			   customer_type, company_name, tax_id, address, birthday,
			   notes, preferences, is_active, created_at, updated_at
		FROM customers 
		WHERE is_active = false
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryWithTenant(ctx, tenantID, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list inactive customers: %w", err)
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		customer := &model.Customer{}
		var email, phone, companyName, taxID, address, notes sql.NullString
		var birthday sql.NullTime

		err := rows.Scan(
			&customer.ID,
			&customer.TenantID,
			&customer.FirstName,
			&customer.LastName,
			&email,
			&phone,
			&customer.CustomerType,
			&companyName,
			&taxID,
			&address,
			&birthday,
			&notes,
			&customer.Preferences,
			&customer.IsActive,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer: %w", err)
		}

		// Convert nullable fields
		customer.Email = StringFromNull(email)
		customer.Phone = StringFromNull(phone)
		customer.CompanyName = StringFromNull(companyName)
		customer.TaxID = StringFromNull(taxID)
		customer.Address = StringFromNull(address)
		customer.Notes = StringFromNull(notes)
		customer.Birthday = TimeFromNull(birthday)

		customers = append(customers, customer)
	}

	return customers, total, nil
}

// Count counts total customers
func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID,
		"SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}

	return count, nil
}

// CountByType counts customers by type
func (r *customerRepository) CountByType(ctx context.Context, customerType string) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID,
		"SELECT COUNT(*) FROM customers WHERE customer_type = $1", customerType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers by type: %w", err)
	}

	return count, nil
}

// CountActive counts active customers
func (r *customerRepository) CountActive(ctx context.Context) (int64, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRowWithTenant(ctx, tenantID,
		"SELECT COUNT(*) FROM customers WHERE is_active = true").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active customers: %w", err)
	}

	return count, nil
}

// ExistsByEmail checks if a customer exists by email
func (r *customerRepository) ExistsByEmail(ctx context.Context, email string, excludeID *int64) (bool, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := "SELECT COUNT(*) FROM customers WHERE email = $1"
	args := []interface{}{email}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	var count int
	err = r.db.QueryRowWithTenant(ctx, tenantID, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// ExistsByTaxID checks if a customer exists by tax ID
func (r *customerRepository) ExistsByTaxID(ctx context.Context, taxID string, excludeID *int64) (bool, error) {
	tenantID, err := GetTenantIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	query := "SELECT COUNT(*) FROM customers WHERE tax_id = $1"
	args := []interface{}{taxID}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	var count int
	err = r.db.QueryRowWithTenant(ctx, tenantID, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check tax ID existence: %w", err)
	}

	return count > 0, nil
}
