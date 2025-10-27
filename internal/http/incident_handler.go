package http

import (
	"encoding/json"
	"net/http"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IncidentHandler struct {
	incidentRepo repository.IncidentRepository
}

func NewIncidentHandler(repo repository.IncidentRepository) *IncidentHandler {
	return &IncidentHandler{incidentRepo: repo}
}

func (h *IncidentHandler) GetIncidentByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	incident, err := h.incidentRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "incident not found")
		return
	}

	respondJSON(w, http.StatusOK, incident)
}

func (h *IncidentHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		respondError(w, http.StatusBadRequest, "status is required")
		return
	}

	status := models.IncidentStatus(statusStr)
	incidents, err := h.incidentRepo.GetByStatus(r.Context(), status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, incidents)
}

type UpdateIncidentDTO struct {
	Status            *string `json:"status,omitempty"`
	AssignedTo        *string `json:"assigned_to,omitempty"`
	RootCause         *string `json:"root_cause,omitempty"`
	CorrectiveActions *string `json:"corrective_actions,omitempty"`
}

func (h *IncidentHandler) UpdateIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var dto UpdateIncidentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	incident, err := h.incidentRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "incident not found")
		return
	}

	if dto.Status != nil {
		incident.Status = models.IncidentStatus(*dto.Status)
	}
	if dto.RootCause != nil {
		incident.RootCause = *dto.RootCause
	}
	if dto.CorrectiveActions != nil {
		incident.CorrectiveActions = *dto.CorrectiveActions
	}
	if dto.AssignedTo != nil {
		assignedTo, err := uuid.Parse(*dto.AssignedTo)
		if err == nil {
			incident.AssignedTo = &assignedTo
		}
	}

	if err := h.incidentRepo.Update(r.Context(), incident); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, incident)
}
