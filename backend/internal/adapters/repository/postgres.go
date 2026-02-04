package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pack-calculator/internal/ports"

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

	var sizes []int
	err := r.db.QueryRowContext(ctx, query).Scan(&sizes)
	if err == sql.ErrNoRows {
		return []int{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active pack sizes: %w", err)
	}

	return sizes, nil
}

func (r *PostgresRepository) Create(sizes []int) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Append-only versioning: deactivate all previous versions and create new one atomically.
	// This ensures only one active version exists at any time while preserving history.
	var maxVersion int
	err = tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM pack_sizes").Scan(&maxVersion)
	if err != nil {
		return fmt.Errorf("failed to get max version: %w", err)
	}

	updateQuery := "UPDATE pack_sizes SET is_active = false WHERE is_active = true"
	_, err = tx.ExecContext(ctx, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to deactivate old versions: %w", err)
	}

	insertQuery := `
		INSERT INTO pack_sizes (version, sizes, is_active) 
		VALUES ($1, $2, true)
	`
	_, err = tx.ExecContext(ctx, insertQuery, maxVersion+1, sizes)
	if err != nil {
		return fmt.Errorf("failed to insert new pack sizes: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

var _ ports.PackSizeRepository = (*PostgresRepository)(nil)
