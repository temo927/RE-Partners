package repository

import (
	"errors"
	"testing"

	pkgerrors "pack-calculator/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresRepository_GetAllActive(t *testing.T) {
	// This test requires a real database connection
	// Skip if running in CI without database
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	dsn := "host=localhost port=5432 user=packcalc password=packcalc dbname=packcalc_test sslmode=disable"
	repo, err := NewPostgresRepository(dsn)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to database: %v", err)
	}
	defer repo.Close()

	t.Run("empty database returns empty slice", func(t *testing.T) {
		sizes, err := repo.GetAllActive()
		if err != nil {
			t.Errorf("GetAllActive() error = %v, want nil", err)
		}
		if sizes == nil {
			t.Error("GetAllActive() returned nil, want empty slice")
		}
		if len(sizes) != 0 {
			t.Errorf("GetAllActive() = %v, want empty slice", sizes)
		}
	})
}

func TestPostgresRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	dsn := "host=localhost port=5432 user=packcalc password=packcalc dbname=packcalc_test sslmode=disable"
	repo, err := NewPostgresRepository(dsn)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to database: %v", err)
	}
	defer repo.Close()

	t.Run("create pack sizes successfully", func(t *testing.T) {
		sizes := []int{250, 500, 1000}
		err := repo.Create(sizes)
		if err != nil {
			t.Errorf("Create() error = %v, want nil", err)
		}

		// Verify it was created
		active, err := repo.GetAllActive()
		if err != nil {
			t.Errorf("GetAllActive() error = %v", err)
		}
		if len(active) != len(sizes) {
			t.Errorf("GetAllActive() = %v, want %v", active, sizes)
		}
	})

	t.Run("create new version deactivates old", func(t *testing.T) {
		oldSizes := []int{250, 500}
		newSizes := []int{100, 200, 300}

		err := repo.Create(oldSizes)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		err = repo.Create(newSizes)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		active, err := repo.GetAllActive()
		if err != nil {
			t.Errorf("GetAllActive() error = %v", err)
		}
		if len(active) != len(newSizes) {
			t.Errorf("GetAllActive() = %v, want %v", active, newSizes)
		}
	})
}

func TestPostgresRepository_ErrorWrapping(t *testing.T) {
	t.Run("invalid DSN returns error", func(t *testing.T) {
		_, err := NewPostgresRepository("invalid dsn")
		if err == nil {
			t.Error("NewPostgresRepository() error = nil, want error")
		}
	})

	t.Run("repository errors are wrapped", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test")
		}

		dsn := "host=invalid port=5432 user=packcalc password=packcalc dbname=packcalc sslmode=disable"
		repo, err := NewPostgresRepository(dsn)
		if err == nil {
			// If connection succeeds, test GetAllActive with closed connection
			repo.Close()
			_, err = repo.GetAllActive()
			if err != nil {
				// Check if error is wrapped with ErrRepository
				if !errors.Is(err, pkgerrors.ErrRepository) {
					t.Errorf("GetAllActive() error should be wrapped with ErrRepository, got: %v", err)
				}
			}
		}
	})
}
