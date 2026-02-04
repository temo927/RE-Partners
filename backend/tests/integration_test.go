//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pack-calculator/internal/app"
	"pack-calculator/internal/domain"
	httptransport "pack-calculator/internal/transport/http"
	pkgerrors "pack-calculator/pkg/errors"
)

type mockPackService struct {
	getPackSizesFunc    func() ([]int, error)
	updatePackSizesFunc func(sizes []int) error
	calculatePacksFunc  func(items int) ([]domain.Pack, error)
}

func (m *mockPackService) GetPackSizes() ([]int, error) {
	if m.getPackSizesFunc != nil {
		return m.getPackSizesFunc()
	}
	return nil, nil
}

func (m *mockPackService) UpdatePackSizes(sizes []int) error {
	if m.updatePackSizesFunc != nil {
		return m.updatePackSizesFunc(sizes)
	}
	return nil
}

func (m *mockPackService) CalculatePacks(items int) ([]domain.Pack, error) {
	if m.calculatePacksFunc != nil {
		return m.calculatePacksFunc(items)
	}
	return nil, nil
}

func setupIntegrationTest(t *testing.T) (*httptransport.Handler, func()) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// In a real integration test, you would set up actual database and cache connections
	// For now, we'll use mocks but test the full HTTP flow
	mockService := &mockPackService{
		getPackSizesFunc: func() ([]int, error) {
			return []int{250, 500, 1000, 2000, 5000}, nil
		},
		updatePackSizesFunc: func(sizes []int) error {
			return nil
		},
		calculatePacksFunc: func(items int) ([]domain.Pack, error) {
			if items <= 0 {
				return nil, pkgerrors.ErrItemsInvalid
			}
			calcService := app.NewCalculationService()
			packSizes := []int{250, 500, 1000, 2000, 5000}
			return calcService.CalculatePacks(packSizes, items), nil
		},
	}

	handler := httptransport.NewHandler(mockService)
	cleanup := func() {}

	return handler, cleanup
}

func TestIntegration_GetPackSizes(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/pack-sizes", nil)
	w := httptest.NewRecorder()

	handler.GetPackSizes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetPackSizes() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response struct {
		Sizes []int `json:"sizes"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("GetPackSizes() invalid JSON: %v", err)
	}

	if len(response.Sizes) == 0 {
		t.Error("GetPackSizes() sizes empty")
	}
}

func TestIntegration_UpdatePackSizes(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	body := map[string]interface{}{
		"sizes": []int{100, 200, 300},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/pack-sizes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdatePackSizes(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("UpdatePackSizes() status = %v, want %v", w.Code, http.StatusNoContent)
	}
}

func TestIntegration_CalculatePacks(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		items          int
		expectedStatus int
	}{
		{
			name:           "small order",
			items:          251,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "large order",
			items:          12001,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "edge case 500000",
			items:          500000,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]interface{}{
				"items": tt.items,
			}
			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/api/calculate", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CalculatePacks(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("CalculatePacks() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if w.Code == http.StatusOK {
				var response struct {
					Packs []struct {
						Size     int `json:"size"`
						Quantity int `json:"quantity"`
					} `json:"packs"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("CalculatePacks() invalid JSON: %v", err)
				}

				if len(response.Packs) == 0 {
					t.Error("CalculatePacks() packs empty")
				}

				// Verify total items >= requested
				totalItems := 0
				for _, pack := range response.Packs {
					totalItems += pack.Size * pack.Quantity
				}
				if totalItems < tt.items {
					t.Errorf("CalculatePacks() total items = %v, want >= %v", totalItems, tt.items)
				}
			}
		})
	}
}

func TestIntegration_HealthCheck(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Health() invalid JSON: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Health() status = %v, want ok", response["status"])
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	handler, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("invalid JSON in request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/pack-sizes", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdatePackSizes(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("UpdatePackSizes() invalid JSON status = %v, want %v", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/calculate", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CalculatePacks(w, req)

		// Should handle missing items field
		if w.Code == http.StatusOK {
			t.Error("CalculatePacks() missing items should return error")
		}
	})
}
