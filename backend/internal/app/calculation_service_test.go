package app

import (
	"reflect"
	"testing"

	"pack-calculator/internal/domain"
)

func TestCalculationService_CalculatePacks(t *testing.T) {
	service := NewCalculationService()

	tests := []struct {
		name      string
		packSizes []int
		items     int
		want      []domain.Pack
	}{
		{
			name:      "single pack exact match",
			packSizes: []int{250},
			items:     250,
			want:      []domain.Pack{{Size: 250, Quantity: 1}},
		},
		{
			name:      "single pack over",
			packSizes: []int{250},
			items:     251,
			want:      []domain.Pack{{Size: 250, Quantity: 2}},
		},
		{
			name:      "multiple packs optimal",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			items:     251,
			want:      []domain.Pack{{Size: 500, Quantity: 1}},
		},
		{
			name:      "multiple packs complex",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			items:     501,
			want:      []domain.Pack{{Size: 500, Quantity: 1}, {Size: 250, Quantity: 1}},
		},
		{
			name:      "large order",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			items:     12001,
			want:      []domain.Pack{{Size: 5000, Quantity: 2}, {Size: 2000, Quantity: 1}, {Size: 250, Quantity: 1}},
		},
		{
			name:      "edge case from requirements",
			packSizes: []int{23, 31, 53},
			items:     500000,
			want:      []domain.Pack{{Size: 23, Quantity: 2}, {Size: 31, Quantity: 7}, {Size: 53, Quantity: 9429}},
		},
		{
			name:      "empty pack sizes",
			packSizes: []int{},
			items:     100,
			want:      []domain.Pack{},
		},
		{
			name:      "zero items",
			packSizes: []int{250, 500},
			items:     0,
			want:      []domain.Pack{},
		},
		{
			name:      "one item",
			packSizes: []int{250, 500},
			items:     1,
			want:      []domain.Pack{{Size: 250, Quantity: 1}},
		},
		{
			name:      "exact match multiple packs",
			packSizes: []int{250, 500},
			items:     750,
			want:      []domain.Pack{{Size: 500, Quantity: 1}, {Size: 250, Quantity: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.CalculatePacks(tt.packSizes, tt.items)

			if len(got) != len(tt.want) {
				t.Errorf("CalculatePacks() returned %d packs, want %d", len(got), len(tt.want))
				return
			}

			gotMap := make(map[int]int)
			for _, p := range got {
				gotMap[p.Size] = p.Quantity
			}

			wantMap := make(map[int]int)
			for _, p := range tt.want {
				wantMap[p.Size] = p.Quantity
			}

			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("CalculatePacks() = %v, want %v", gotMap, wantMap)
			}

			if len(got) > 0 {
				totalItems := 0
				for _, p := range got {
					totalItems += p.Size * p.Quantity
				}

				if totalItems < tt.items {
					t.Errorf("Total items %d is less than required %d", totalItems, tt.items)
				}
			}
		})
	}
}

func TestCalculationService_EdgeCase(t *testing.T) {
	service := NewCalculationService()
	packSizes := []int{23, 31, 53}
	items := 500000

	result := service.CalculatePacks(packSizes, items)

	expected := map[int]int{23: 2, 31: 7, 53: 9429}
	resultMap := make(map[int]int)
	for _, p := range result {
		resultMap[p.Size] = p.Quantity
	}

	if !reflect.DeepEqual(resultMap, expected) {
		t.Errorf("Edge case failed: got %v, want %v", resultMap, expected)
	}

	total := 0
	for _, p := range result {
		total += p.Size * p.Quantity
	}

	if total != 500000 {
		t.Errorf("Total items %d != 500000", total)
	}
}

func TestCalculationService_MinimizeItems(t *testing.T) {
	service := NewCalculationService()

	tests := []struct {
		name      string
		packSizes []int
		items     int
		check     func([]domain.Pack) bool
	}{
		{
			name:      "prefer larger pack to minimize items",
			packSizes: []int{250, 500},
			items:     251,
			check: func(packs []domain.Pack) bool {
				total := 0
				for _, p := range packs {
					total += p.Size * p.Quantity
				}
				return total == 500
			},
		},
		{
			name:      "minimize items over pack count",
			packSizes: []int{250, 500, 1000},
			items:     501,
			check: func(packs []domain.Pack) bool {
				total := 0
				for _, p := range packs {
					total += p.Size * p.Quantity
				}
				return total == 750
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculatePacks(tt.packSizes, tt.items)
			if !tt.check(result) {
				t.Errorf("CalculatePacks() did not minimize items correctly")
			}
		})
	}
}
