package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pack-calculator/internal/domain"
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

func TestHandler_GetPackSizes(t *testing.T) {
	tests := []struct {
		name           string
		mockService    *mockPackService
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			mockService: &mockPackService{
				getPackSizesFunc: func() ([]int, error) {
					return []int{250, 500, 1000}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"sizes": []interface{}{250.0, 500.0, 1000.0},
			},
		},
		{
			name: "repository error",
			mockService: &mockPackService{
				getPackSizesFunc: func() ([]int, error) {
					return nil, pkgerrors.ErrRepository
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mockService)
			req := httptest.NewRequest("GET", "/api/pack-sizes", nil)
			w := httptest.NewRecorder()

			handler.GetPackSizes(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetPackSizes() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != nil {
				var got map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("GetPackSizes() invalid JSON response: %v", err)
				}
			}
		})
	}
}

func TestHandler_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockService    *mockPackService
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]interface{}{
				"sizes": []int{250, 500, 1000},
			},
			mockService: &mockPackService{
				updatePackSizesFunc: func(sizes []int) error {
					return nil
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "empty sizes",
			body: map[string]interface{}{
				"sizes": []int{},
			},
			mockService: &mockPackService{
				updatePackSizesFunc: func(sizes []int) error {
					return pkgerrors.ErrPackSizesEmpty
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid JSON",
			body: "invalid",
			mockService: &mockPackService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repository error",
			body: map[string]interface{}{
				"sizes": []int{250, 500},
			},
			mockService: &mockPackService{
				updatePackSizesFunc: func(sizes []int) error {
					return pkgerrors.ErrRepository
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			handler := NewHandler(tt.mockService)
			req := httptest.NewRequest("POST", "/api/pack-sizes", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdatePackSizes(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("UpdatePackSizes() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestHandler_CalculatePacks(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		mockService    *mockPackService
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]interface{}{
				"items": 251,
			},
			mockService: &mockPackService{
				calculatePacksFunc: func(items int) ([]domain.Pack, error) {
					return []domain.Pack{
						{Size: 250, Quantity: 1},
						{Size: 1, Quantity: 1},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid items",
			body: map[string]interface{}{
				"items": 0,
			},
			mockService: &mockPackService{
				calculatePacksFunc: func(items int) ([]domain.Pack, error) {
					return nil, pkgerrors.ErrItemsInvalid
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid JSON",
			body: "invalid",
			mockService: &mockPackService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: map[string]interface{}{
				"items": 100,
			},
			mockService: &mockPackService{
				calculatePacksFunc: func(items int) ([]domain.Pack, error) {
					return nil, pkgerrors.ErrRepository
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			handler := NewHandler(tt.mockService)
			req := httptest.NewRequest("POST", "/api/calculate", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CalculatePacks(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("CalculatePacks() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestHandler_Health(t *testing.T) {
	handler := NewHandler(&mockPackService{})
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Health() invalid JSON response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Health() status = %v, want ok", response["status"])
	}
}

func TestHandler_handleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "not found",
			err:            pkgerrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid input",
			err:            pkgerrors.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "pack sizes empty",
			err:            pkgerrors.ErrPackSizesEmpty,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "items invalid",
			err:            pkgerrors.ErrItemsInvalid,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "repository error",
			err:            pkgerrors.ErrRepository,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "cache error",
			err:            pkgerrors.ErrCache,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unknown error",
			err:            pkgerrors.Wrap(pkgerrors.ErrRepository, "test"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&mockPackService{})
			w := httptest.NewRecorder()

			handler.handleError(w, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleError() status = %v, want %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

