package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/config"
	_ "github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Healthcheck verifies the database connection is healthy
func (db *DB) Healthcheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

// Context keys for tenant information
type contextKey string

const (
	TenantIDKey contextKey = "tenant_id"
)

// WithTenantID adds tenant ID (UUID string) to context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID extracts tenant ID from context
func GetTenantID(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	return tenantID, ok
}

// SetTenantID sets the tenant ID in the database session for RLS
func (db *DB) SetTenantID(ctx context.Context, tenantID string) error {
	// PostgreSQL SET command doesn't accept placeholders, must use string formatting
	// Safe because tenant_id is validated as UUID format
	query := fmt.Sprintf("SET app.current_tenant_id = '%s'", tenantID)
	_, err := db.ExecContext(ctx, query)
	return err
}

// BeginTx starts a new transaction with tenant ID set
func (db *DB) BeginTxWithTenant(ctx context.Context, tenantID string) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Set tenant ID for RLS
	// PostgreSQL SET command doesn't accept placeholders
	query := fmt.Sprintf("SET app.current_tenant_id = '%s'", tenantID)
	if _, err := tx.ExecContext(ctx, query); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set tenant ID: %w", err)
	}

	return tx, nil
}

// ExecWithTenant executes a query with tenant ID set
func (db *DB) ExecWithTenant(ctx context.Context, tenantID string, query string, args ...interface{}) (sql.Result, error) {
	// Set tenant ID
	if err := db.SetTenantID(ctx, tenantID); err != nil {
		return nil, err
	}

	return db.ExecContext(ctx, query, args...)
}

// QueryWithTenant executes a query with tenant ID set
func (db *DB) QueryWithTenant(ctx context.Context, tenantID string, query string, args ...interface{}) (*sql.Rows, error) {
	// Set tenant ID
	if err := db.SetTenantID(ctx, tenantID); err != nil {
		return nil, err
	}

	return db.QueryContext(ctx, query, args...)
}

// QueryRowWithTenant executes a query that returns a single row with tenant ID set
func (db *DB) QueryRowWithTenant(ctx context.Context, tenantID string, query string, args ...interface{}) *sql.Row {
	// Set tenant ID - ignore error as we can't return it from QueryRow
	db.SetTenantID(ctx, tenantID)

	return db.QueryRowContext(ctx, query, args...)
}

// GetTenantIDFromContext extracts tenant ID from context
func GetTenantIDFromContext(ctx context.Context) (string, error) {
	tenantID, ok := GetTenantID(ctx)
	if !ok {
		return "", fmt.Errorf("tenant ID not found in context")
	}
	return tenantID, nil
}

// WithTenantContext creates a context with tenant ID
func WithTenantContext(ctx context.Context, tenantID string) context.Context {
	return WithTenantID(ctx, tenantID)
}

// Transaction helper function that sets tenant ID and runs a function in a transaction
func (db *DB) TransactionWithTenant(ctx context.Context, tenantID string, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTxWithTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ScanRowsToMap scans SQL rows into a map slice (utility function)
func ScanRowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}

			rowMap[col] = v
		}

		results = append(results, rowMap)
	}

	return results, nil
}

// NullString helper for handling nullable strings
func NullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// StringFromNull helper for converting nullable strings
func StringFromNull(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// NullTime helper for handling nullable times
func NullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// TimeFromNull helper for converting nullable times
func TimeFromNull(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// NullInt64 helper for handling nullable int64
func NullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// Int64FromNull helper for converting nullable int64
func Int64FromNull(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}
