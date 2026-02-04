package app

import (
	"fmt"

	"pack-calculator/internal/domain"
	"pack-calculator/internal/ports"
)

type PackService struct {
	repo  ports.PackSizeRepository
	cache ports.Cache
}

func NewPackService(repo ports.PackSizeRepository, cache ports.Cache) *PackService {
	return &PackService{
		repo:  repo,
		cache: cache,
	}
}

func (s *PackService) GetPackSizes() ([]int, error) {
	cacheKey := "pack-sizes:active"

	sizes, err := s.cache.Get(cacheKey)
	if err == nil {
		return sizes, nil
	}

	sizes, err = s.repo.GetAllActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes from repository: %w", err)
	}

	if err := s.cache.Set(cacheKey, sizes, 3600); err != nil {
		// Log cache error but don't fail the request
		_ = err
	}

	return sizes, nil
}

func (s *PackService) UpdatePackSizes(sizes []int) error {
	if err := s.repo.Create(sizes); err != nil {
		return fmt.Errorf("failed to create pack sizes: %w", err)
	}

	// Invalidate cache by deleting the active key
	// New requests will fetch from DB and cache with new version
	cacheKey := "pack-sizes:active"
	if err := s.cache.Delete(cacheKey); err != nil {
		// Log cache error but don't fail the request
		_ = err
	}

	return nil
}

func (s *PackService) CalculatePacks(items int) ([]domain.Pack, error) {
	packSizes, err := s.GetPackSizes()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}

	calcService := NewCalculationService()
	return calcService.CalculatePacks(packSizes, items), nil
}
