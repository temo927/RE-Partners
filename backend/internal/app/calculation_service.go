package app

import "pack-calculator/internal/domain"

type CalculationService struct{}

func NewCalculationService() *CalculationService {
	return &CalculationService{}
}

func (s *CalculationService) CalculatePacks(packSizes []int, items int) []domain.Pack {
	if len(packSizes) == 0 || items <= 0 {
		return []domain.Pack{}
	}

	maxSize := 0
	for _, size := range packSizes {
		if size > maxSize {
			maxSize = size
		}
	}

	// We search up to items + maxSize because we can only use whole packs.
	// If exact match isn't possible, we may need to send more items than requested.
	maxTarget := items + maxSize
	minItems, minPacks := s.findOptimalSolution(packSizes, items, maxTarget)

	if minItems == -1 {
		return []domain.Pack{}
	}

	resultMap := s.reconstructSolution(packSizes, items, maxTarget, minItems, minPacks)
	return s.mapToPacks(resultMap)
}

func (s *CalculationService) mapToPacks(resultMap map[int]int) []domain.Pack {
	var packs []domain.Pack
	for size, quantity := range resultMap {
		packs = append(packs, domain.Pack{
			Size:     size,
			Quantity: quantity,
		})
	}
	return packs
}

type solution struct {
	totalItems int
	packCount  int
}

// 1. Minimizes total items sent (primary objective)
// 2. Minimizes number of packs (secondary objective, when items are equal)
func (s *CalculationService) findOptimalSolution(packSizes []int, items int, maxTarget int) (int, int) {
	dp := make(map[int]solution)
	dp[0] = solution{totalItems: 0, packCount: 0}

	for target := 1; target <= maxTarget; target++ {
		best := solution{totalItems: maxTarget + 1, packCount: maxTarget + 1}

		for _, size := range packSizes {
			if size <= target {
				prev := target - size
				if prevSol, exists := dp[prev]; exists {
					candidate := solution{
						totalItems: prevSol.totalItems + size,
						packCount:  prevSol.packCount + 1,
					}

					if candidate.totalItems < best.totalItems ||
						(candidate.totalItems == best.totalItems && candidate.packCount < best.packCount) {
						best = candidate
					}
				}
			}
		}

		if best.totalItems <= maxTarget {
			dp[target] = best
		}
	}

	minItems := maxTarget + 1
	minPacks := maxTarget + 1

	for target := items; target <= maxTarget; target++ {
		if sol, exists := dp[target]; exists {
			if sol.totalItems < minItems ||
				(sol.totalItems == minItems && sol.packCount < minPacks) {
				minItems = sol.totalItems
				minPacks = sol.packCount
			}
		}
	}

	if minItems > maxTarget {
		return -1, -1
	}

	return minItems, minPacks
}

func (s *CalculationService) reconstructSolution(packSizes []int, items int, maxTarget int, minItems int, minPacks int) map[int]int {
	result := make(map[int]int)

	var backtrack func(target int, current map[int]int, totalItems int, packCount int) bool
	backtrack = func(target int, current map[int]int, totalItems int, packCount int) bool {
		if target >= items && totalItems == minItems && packCount == minPacks {
			for k, v := range current {
				result[k] = v
			}
			return true
		}

		if target > maxTarget || totalItems > minItems || packCount > minPacks {
			return false
		}

		for _, size := range packSizes {
			if size <= target {
				current[size]++
				if backtrack(target-size, current, totalItems+size, packCount+1) {
					return true
				}
				current[size]--
				if current[size] == 0 {
					delete(current, size)
				}
			}
		}

		return false
	}

	for target := items; target <= maxTarget; target++ {
		if backtrack(target, make(map[int]int), 0, 0) {
			break
		}
	}

	return result
}
