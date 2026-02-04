package http

import (
	"encoding/json"
	"net/http"

	"pack-calculator/internal/app"
	"pack-calculator/internal/domain"
	"pack-calculator/internal/transport"
)

type Handler struct {
	packService *app.PackService
}

func NewHandler(packService *app.PackService) *Handler {
	return &Handler{
		packService: packService,
	}
}

func (h *Handler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.packService.GetPackSizes()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := transport.PackSizesResponse{Sizes: sizes}
	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	var req transport.UpdatePackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	if len(req.Sizes) == 0 {
		h.writeError(w, http.StatusBadRequest, &validationError{message: "sizes cannot be empty"})
		return
	}

	if err := h.packService.UpdatePackSizes(req.Sizes); err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	var req transport.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Items <= 0 {
		h.writeError(w, http.StatusBadRequest, &validationError{message: "items must be greater than 0"})
		return
	}

	packs, err := h.packService.CalculatePacks(req.Items)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := transport.CalculateResponse{
		Packs: h.domainPacksToResponse(packs),
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) domainPacksToResponse(packs []domain.Pack) []transport.PackResponse {
	result := make([]transport.PackResponse, len(packs))
	for i, p := range packs {
		result[i] = transport.PackResponse{
			Size:     p.Size,
			Quantity: p.Quantity,
		}
	}
	return result
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(transport.ErrorResponse{Error: err.Error()})
}

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}
