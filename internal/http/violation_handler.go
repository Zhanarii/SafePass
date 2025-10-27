package http

import (
	"encoding/json"
	"net/http"
	"time"

	"ppe-detection/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ViolationHandler struct {
	violationSvc *service.ViolationService
}

func NewViolationHandler(svc *service.ViolationService) *ViolationHandler {
	return &ViolationHandler{violationSvc: svc}
}

func (h *ViolationHandler) GetViolationByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	violation, err := h.violationSvc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "violation not found")
		return
	}

	respondJSON(w, http.StatusOK, violation)
}

func (h *ViolationHandler) AcknowledgeViolation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var dto struct {
		AcknowledgedBy string `json:"acknowledged_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID, err := uuid.Parse(dto.AcknowledgedBy)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	if err := h.violationSvc.AcknowledgeViolation(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "acknowledged"})
}

func (h *ViolationHandler) CreateIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	violationID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid violation_id")
		return
	}

	incident, err := h.violationSvc.CreateIncidentFromViolation(r.Context(), violationID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, incident)
}

func (h *ViolationHandler) GetViolationStats(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid from date")
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid to date")
		return
	}

	stats, err := h.violationSvc.GetViolationStats(r.Context(), from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
