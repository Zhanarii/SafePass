package http

import (
	"encoding/json"
	"net/http"
	"time"

	"ppe-detection/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AccessHandler struct {
	accessSvc *service.AccessControlService
}

func NewAccessHandler(svc *service.AccessControlService) *AccessHandler {
	return &AccessHandler{accessSvc: svc}
}

type CheckAccessDTO struct {
	UserID       *string `json:"user_id,omitempty"`
	BadgeNumber  *string `json:"badge_number,omitempty"`
	CameraID     string  `json:"camera_id"`
	AccessZoneID string  `json:"access_zone_id"`
	DetectionID  *string `json:"detection_id,omitempty"`
}

func (h *AccessHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	var dto CheckAccessDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cameraID, err := uuid.Parse(dto.CameraID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid camera_id")
		return
	}

	zoneID, err := uuid.Parse(dto.AccessZoneID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid access_zone_id")
		return
	}

	req := &service.AccessRequest{
		CameraID:     cameraID,
		AccessZoneID: zoneID,
		BadgeNumber:  dto.BadgeNumber,
	}

	if dto.UserID != nil {
		userID, err := uuid.Parse(*dto.UserID)
		if err == nil {
			req.UserID = &userID
		}
	}

	if dto.DetectionID != nil {
		detectionID, err := uuid.Parse(*dto.DetectionID)
		if err == nil {
			req.DetectionID = &detectionID
		}
	}

	resp, err := h.accessSvc.CheckAccess(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

type OverrideAccessDTO struct {
	OverrideBy string `json:"override_by"`
	Reason     string `json:"reason"`
}

func (h *AccessHandler) OverrideAccess(w http.ResponseWriter, r *http.Request) {
	logIDStr := chi.URLParam(r, "log_id")
	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid log_id")
		return
	}

	var dto OverrideAccessDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	overrideBy, err := uuid.Parse(dto.OverrideBy)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid override_by")
		return
	}

	if err := h.accessSvc.OverrideAccess(r.Context(), logID, overrideBy, dto.Reason); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "access_granted"})
}

func (h *AccessHandler) GetAccessHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

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

	history, err := h.accessSvc.GetAccessHistory(r.Context(), userID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, history)
}
