package app

import (
	"errors"
	"reflect"
	"testing"

	"pack-calculator/internal/domain"
	"pack-calculator/internal/ports"
)

type mockRepository struct {
	getAllActiveFunc func() ([]int, error)
	createFunc       func(sizes []int) error
}

func (m *mockRepository) GetAllActive() ([]int, error) {
	if m.getAllActiveFunc != nil {
		return m.getAllActiveFunc()
	}
	return nil, nil
}

func (m *mockRepository) Create(sizes []int) error {
	if m.createFunc != nil {
		return m.createFunc(sizes)
	}
	return nil
}

type mockCache struct {
	getFunc    func(key string) ([]int, error)
	setFunc    func(key string, value []int, ttl int) error
	deleteFunc func(key string) error
}

func (m *mockCache) Get(key string) ([]int, error) {
	if m.getFunc != nil {
		return m.getFunc(key)
	}
	return nil, errors.New("key not found")
}

func (m *mockCache) Set(key string, value []int, ttl int) error {
	if m.setFunc != nil {
		return m.setFunc(key, value, ttl)
	}
	return nil
}

func (m *mockCache) Delete(key string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(key)
	}
	return nil
}

func TestPackService_GetPackSizes(t *testing.T) {
	tests := []struct {
		name           string
		repo           ports.PackSizeRepository
		cache          ports.Cache
		want           []int
		wantErr        bool
		cacheCalled    bool
		repoCalled     bool
		cacheSetCalled bool
	}{
		{
			name: "cache hit",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return []int{250, 500, 1000}, nil
				},
			},
			want:        []int{250, 500, 1000},
			wantErr:     false,
			cacheCalled: true,
			repoCalled:  false,
		},
		{
			name: "cache miss, repository success",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return nil, errors.New("key not found")
				},
				setFunc: func(key string, value []int, ttl int) error {
					return nil
				},
			},
			repo: &mockRepository{
				getAllActiveFunc: func() ([]int, error) {
					return []int{250, 500, 1000}, nil
				},
			},
			want:           []int{250, 500, 1000},
			wantErr:        false,
			cacheCalled:    true,
			repoCalled:     true,
			cacheSetCalled: true,
		},
		{
			name: "cache miss, repository error",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return nil, errors.New("key not found")
				},
			},
			repo: &mockRepository{
				getAllActiveFunc: func() ([]int, error) {
					return nil, errors.New("database error")
				},
			},
			want:        nil,
			wantErr:     true,
			cacheCalled: true,
			repoCalled:  true,
		},
		{
			name: "cache set error doesn't fail request",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return nil, errors.New("key not found")
				},
				setFunc: func(key string, value []int, ttl int) error {
					return errors.New("cache set failed")
				},
			},
			repo: &mockRepository{
				getAllActiveFunc: func() ([]int, error) {
					return []int{250, 500}, nil
				},
			},
			want:           []int{250, 500},
			wantErr:        false,
			cacheCalled:    true,
			repoCalled:     true,
			cacheSetCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcService := NewCalculationService()
			service := NewPackService(tt.repo, tt.cache, calcService)
			got, err := service.GetPackSizes()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackSizes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPackSizes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackService_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name         string
		repo         ports.PackSizeRepository
		cache        ports.Cache
		sizes        []int
		wantErr      bool
		repoCalled   bool
		cacheCalled  bool
		deleteCalled bool
	}{
		{
			name: "successful update",
			repo: &mockRepository{
				createFunc: func(sizes []int) error {
					return nil
				},
			},
			cache: &mockCache{
				deleteFunc: func(key string) error {
					return nil
				},
			},
			sizes:        []int{250, 500, 1000},
			wantErr:      false,
			repoCalled:   true,
			cacheCalled:  true,
			deleteCalled: true,
		},
		{
			name: "repository error",
			repo: &mockRepository{
				createFunc: func(sizes []int) error {
					return errors.New("database error")
				},
			},
			cache:      &mockCache{},
			sizes:      []int{250, 500},
			wantErr:    true,
			repoCalled: true,
		},
		{
			name: "cache delete error doesn't fail request",
			repo: &mockRepository{
				createFunc: func(sizes []int) error {
					return nil
				},
			},
			cache: &mockCache{
				deleteFunc: func(key string) error {
					return errors.New("cache delete failed")
				},
			},
			sizes:        []int{250, 500},
			wantErr:      false,
			repoCalled:   true,
			cacheCalled:  true,
			deleteCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcService := NewCalculationService()
			service := NewPackService(tt.repo, tt.cache, calcService)
			err := service.UpdatePackSizes(tt.sizes)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackSizes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackService_CalculatePacks(t *testing.T) {
	tests := []struct {
		name    string
		repo    ports.PackSizeRepository
		cache   ports.Cache
		items   int
		want    []domain.Pack
		wantErr bool
	}{
		{
			name: "successful calculation",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return []int{250, 500, 1000}, nil
				},
			},
			items:   251,
			want:    []domain.Pack{{Size: 500, Quantity: 1}},
			wantErr: false,
		},
		{
			name: "get pack sizes error",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return nil, errors.New("key not found")
				},
			},
			repo: &mockRepository{
				getAllActiveFunc: func() ([]int, error) {
					return nil, errors.New("database error")
				},
			},
			items:   100,
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty pack sizes",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return []int{}, nil
				},
			},
			items:   100,
			want:    []domain.Pack{},
			wantErr: false,
		},
		{
			name: "zero items",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return []int{250, 500}, nil
				},
			},
			items:   0,
			want:    nil,
			wantErr: true,
		},
		{
			name: "items out of range (too large)",
			cache: &mockCache{
				getFunc: func(key string) ([]int, error) {
					return []int{250, 500}, nil
				},
			},
			items:   2147483648,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcService := NewCalculationService()
			service := NewPackService(tt.repo, tt.cache, calcService)
			got, err := service.CalculatePacks(tt.items)

			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePacks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePacks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackService_UpdatePackSizes_Validation(t *testing.T) {
	tests := []struct {
		name    string
		sizes   []int
		wantErr bool
	}{
		{
			name:    "empty sizes",
			sizes:   []int{},
			wantErr: true,
		},
		{
			name:    "pack size too large",
			sizes:   []int{250, 500, 2147483648},
			wantErr: true,
		},
		{
			name:    "pack size negative",
			sizes:   []int{250, -100},
			wantErr: true,
		},
		{
			name:    "pack size zero",
			sizes:   []int{250, 0},
			wantErr: true,
		},
		{
			name:    "duplicate pack sizes",
			sizes:   []int{250, 500, 250},
			wantErr: true,
		},
		{
			name:    "valid sizes at max boundary",
			sizes:   []int{250, 500, 2147483647},
			wantErr: false,
		},
		{
			name:    "valid sizes",
			sizes:   []int{250, 500, 1000},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				createFunc: func(sizes []int) error {
					return nil
				},
			}
			cache := &mockCache{
				deleteFunc: func(key string) error {
					return nil
				},
			}
			calcService := NewCalculationService()
			service := NewPackService(repo, cache, calcService)
			err := service.UpdatePackSizes(tt.sizes)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePackSizes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
