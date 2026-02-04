package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"pack-calculator/internal/ports"
	pkgerrors "pack-calculator/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) GetAllActive() ([]int, error) {
	ctx := context.Background()
	query := `
		SELECT sizes 
		FROM pack_sizes 
		WHERE is_active = true 
		ORDER BY version DESC 
		LIMIT 1
	`

	var arrayStr string
	err := r.db.QueryRowContext(ctx, query).Scan(&arrayStr)
	if err == sql.ErrNoRows {
		return []int{}, nil
	}
	if err != nil {
		return nil, pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to get active pack sizes")
	}

	// Parse PostgreSQL array format: {1,2,3} or {1, 2, 3}
	arrayStr = strings.Trim(arrayStr, "{}")
	if arrayStr == "" {
		return []int{}, nil
	}

	parts := strings.Split(arrayStr, ",")
	sizes := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to parse pack size")
		}
		sizes = append(sizes, val)
	}

	return sizes, nil
}

func (r *PostgresRepository) Create(sizes []int) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Append-only versioning: deactivate all previous versions and create new one atomically.
	// This ensures only one active version exists at any time while preserving history.
	var maxVersion int
	err = tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM pack_sizes").Scan(&maxVersion)
	if err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to get max version")
	}

	updateQuery := "UPDATE pack_sizes SET is_active = false WHERE is_active = true"
	_, err = tx.ExecContext(ctx, updateQuery)
	if err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to deactivate old versions")
	}

	insertQuery := `
		INSERT INTO pack_sizes (version, sizes, is_active) 
		VALUES ($1, $2::integer[], true)
	`
	
	// Format as PostgreSQL array: {1,2,3}
	arrayParts := make([]string, len(sizes))
	for i, size := range sizes {
		arrayParts[i] = strconv.Itoa(size)
	}
	arrayStr := "{" + strings.Join(arrayParts, ",") + "}"
	
	_, err = tx.ExecContext(ctx, insertQuery, maxVersion+1, arrayStr)
	if err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to insert new pack sizes")
	}

	if err := tx.Commit(); err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrRepository, "failed to commit transaction")
	}

	return nil
}

var _ ports.PackSizeRepository = (*PostgresRepository)(nil)
