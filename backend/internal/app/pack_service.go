package app

import (
	"errors"
	"log/slog"

	"pack-calculator/internal/domain"
	"pack-calculator/internal/ports"
	pkgerrors "pack-calculator/pkg/errors"
	"pack-calculator/pkg/logger"
)

type PackServiceInterface interface {
	GetPackSizes() ([]int, error)
	UpdatePackSizes(sizes []int) error
	CalculatePacks(items int) ([]domain.Pack, error)
}

type PackService struct {
	repo           ports.PackSizeRepository
	cache          ports.Cache
	calculationSvc *CalculationService
	logger         *slog.Logger
}

func NewPackService(repo ports.PackSizeRepository, cache ports.Cache, calculationSvc *CalculationService) *PackService {
	return &PackService{
		repo:           repo,
		cache:          cache,
		calculationSvc: calculationSvc,
		logger:         logger.Default(),
	}
}

func (s *PackService) GetPackSizes() ([]int, error) {
	cacheKey := "pack-sizes:active"

	sizes, err := s.cache.Get(cacheKey)
	if err == nil {
		return sizes, nil
	}

	if !errors.Is(err, pkgerrors.ErrNotFound) {
		s.logger.Warn("Cache get failed, falling back to repository", "error", err, "key", cacheKey)
	}

	sizes, err = s.repo.GetAllActive()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get pack sizes from repository")
	}

	if err := s.cache.Set(cacheKey, sizes, 3600); err != nil {
		s.logger.Warn("Failed to set cache", "error", err, "key", cacheKey)
	}

	return sizes, nil
}

func (s *PackService) UpdatePackSizes(sizes []int) error {
	if len(sizes) == 0 {
		return pkgerrors.ErrPackSizesEmpty
	}

	seen := make(map[int]bool)
	for _, size := range sizes {
		if size < pkgerrors.MinPackSize || size > pkgerrors.MaxPackSize {
			return pkgerrors.ErrPackSizeOutOfRange
		}
		if seen[size] {
			return pkgerrors.ErrDuplicatePackSizes
		}
		seen[size] = true
	}

	if err := s.repo.Create(sizes); err != nil {
		return pkgerrors.Wrap(err, "failed to create pack sizes")
	}

	// Invalidate cache by deleting the active key
	// New requests will fetch from DB and cache with new version
	cacheKey := "pack-sizes:active"
	if err := s.cache.Delete(cacheKey); err != nil {
		s.logger.Warn("Failed to delete cache", "error", err, "key", cacheKey)
	}

	return nil
}

func (s *PackService) CalculatePacks(items int) ([]domain.Pack, error) {
	if items < pkgerrors.MinItems || items > pkgerrors.MaxItems {
		return nil, pkgerrors.ErrItemsOutOfRange
	}

	packSizes, err := s.GetPackSizes()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get pack sizes")
	}

	return s.calculationSvc.CalculatePacks(packSizes, items), nil
}
