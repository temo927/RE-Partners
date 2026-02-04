package app

import (
	"sort"

	"pack-calculator/internal/domain"
)

type CalculationService struct{}

func NewCalculationService() *CalculationService {
	return &CalculationService{}
}

func (s *CalculationService) CalculatePacks(packSizes []int, items int) []domain.Pack {
	if len(packSizes) == 0 || items <= 0 {
		return []domain.Pack{}
	}

	maxSize := maxSizeInSlice(packSizes)

	// We search up to items + maxSize because we can only use whole packs.
	// If exact match isn't possible, we may need to send more items than requested.
	result := s.findOptimalCombination(packSizes, items, maxSize)
	return s.mapToPacks(result)
}

func (s *CalculationService) mapToPacks(resultMap map[int]int) []domain.Pack {
	var packs []domain.Pack
	for size, quantity := range resultMap {
		packs = append(packs, domain.Pack{
			Size:     size,
			Quantity: quantity,
		})
	}
	
	// Sort by pack size in ascending order
	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Size < packs[j].Size
	})
	
	return packs
}

type dpState struct {
	totalItems  int
	packCount   int
	combination map[int]int
}

// 1. Minimizes total items sent (primary objective)
// 2. Minimizes number of packs (secondary objective, when items are equal)
func (s *CalculationService) findOptimalCombination(packSizes []int, items int, maxSize int) map[int]int {
	maxTarget := items + maxSize

	// For very large inputs, use optimized approach
	if items > 100000 {
		return s.findOptimalLargeInput(packSizes, items, maxSize)
	}

	dp := make(map[int]*dpState)
	dp[0] = &dpState{totalItems: 0, packCount: 0, combination: make(map[int]int)}

	for target := 1; target <= maxTarget; target++ {
		best := &dpState{totalItems: maxTarget + 1, packCount: maxTarget + 1, combination: nil}

		for _, size := range packSizes {
			if size <= target {
				prev := target - size
				if prevState, exists := dp[prev]; exists {
					newCombination := make(map[int]int)
					for k, v := range prevState.combination {
						newCombination[k] = v
					}
					newCombination[size]++

					candidate := &dpState{
						totalItems:  prevState.totalItems + size,
						packCount:   prevState.packCount + 1,
						combination: newCombination,
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

	bestState := &dpState{totalItems: maxTarget + 1, packCount: maxTarget + 1, combination: nil}
	for target := items; target <= maxTarget; target++ {
		if state, exists := dp[target]; exists {
			if state.totalItems < bestState.totalItems ||
				(state.totalItems == bestState.totalItems && state.packCount < bestState.packCount) {
				bestState = state
			}
		}
	}

	if bestState.combination == nil {
		return make(map[int]int)
	}

	result := make(map[int]int)
	for k, v := range bestState.combination {
		result[k] = v
	}
	return result
}

func (s *CalculationService) findOptimalLargeInput(packSizes []int, items int, maxSize int) map[int]int {
	// Use BFS-like approach for large inputs
	type state struct {
		total     int
		packCount int
		combo     map[int]int
	}

	visited := make(map[int]*state)
	queue := []*state{{total: 0, packCount: 0, combo: make(map[int]int)}}
	visited[0] = queue[0]

	bestState := &state{total: items + maxSize + 1, packCount: items + maxSize + 1, combo: nil}
	maxTarget := items + maxSize

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.total >= items {
			if current.total < bestState.total ||
				(current.total == bestState.total && current.packCount < bestState.packCount) {
				bestState = current
			}
		}

		if current.total >= maxTarget {
			continue
		}

		for _, size := range packSizes {
			nextTotal := current.total + size
			if nextTotal > maxTarget {
				continue
			}

			existing, exists := visited[nextTotal]
			if !exists {
				newCombo := make(map[int]int)
				for k, v := range current.combo {
					newCombo[k] = v
				}
				newCombo[size]++

				newState := &state{
					total:     nextTotal,
					packCount: current.packCount + 1,
					combo:     newCombo,
				}

				visited[nextTotal] = newState
				queue = append(queue, newState)
			} else {
				newPackCount := current.packCount + 1
				if nextTotal < existing.total ||
					(nextTotal == existing.total && newPackCount < existing.packCount) {
					newCombo := make(map[int]int)
					for k, v := range current.combo {
						newCombo[k] = v
					}
					newCombo[size]++
					existing.total = nextTotal
					existing.packCount = newPackCount
					existing.combo = newCombo
				}
			}
		}
	}

	if bestState.combo == nil {
		return make(map[int]int)
	}

	return bestState.combo
}

func maxSizeInSlice(packSizes []int) int {
	max := 0
	for _, size := range packSizes {
		if size > max {
			max = size
		}
	}
	return max
}
