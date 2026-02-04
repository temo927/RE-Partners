package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"pack-calculator/internal/app"
	"pack-calculator/internal/domain"
	"pack-calculator/internal/transport"
	pkgerrors "pack-calculator/pkg/errors"
)

type Handler struct {
	packService app.PackServiceInterface
}

func NewHandler(packService app.PackServiceInterface) *Handler {
	return &Handler{
		packService: packService,
	}
}

func (h *Handler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.packService.GetPackSizes()
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := transport.PackSizesResponse{Sizes: sizes}
	h.writeJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	var req transport.UpdatePackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, pkgerrors.ErrInvalidInput)
		return
	}

	if err := h.packService.UpdatePackSizes(req.Sizes); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	var req transport.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, pkgerrors.ErrInvalidInput)
		return
	}

	packs, err := h.packService.CalculatePacks(req.Items)
	if err != nil {
		h.handleError(w, err)
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

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var status int

	switch {
	case errors.Is(err, pkgerrors.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, pkgerrors.ErrInvalidInput) || errors.Is(err, pkgerrors.ErrPackSizesEmpty) || errors.Is(err, pkgerrors.ErrItemsInvalid) || errors.Is(err, pkgerrors.ErrPackSizeOutOfRange) || errors.Is(err, pkgerrors.ErrItemsOutOfRange) || errors.Is(err, pkgerrors.ErrDuplicatePackSizes):
		status = http.StatusBadRequest
	case errors.Is(err, pkgerrors.ErrRepository) || errors.Is(err, pkgerrors.ErrCache):
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}

	h.writeError(w, status, err)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(transport.ErrorResponse{Error: err.Error()})
}
