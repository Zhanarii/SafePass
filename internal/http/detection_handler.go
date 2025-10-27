// internal/handler/http/detection_handler.go
package http

import (
	"encoding/json"
	"net/http"
	"time"

	"ppe-detection/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DetectionHandler struct {
	detectionSvc *service.DetectionService
}

func NewDetectionHandler(svc *service.DetectionService) *DetectionHandler {
	return &DetectionHandler{detectionSvc: svc}
}

type ProcessFrameDTO struct {
	CameraID         string                 `json:"camera_id"`
	FrameURL         string                 `json:"frame_url"`
	DetectedPPE      []string               `json:"detected_ppe"`
	ConfidenceScores map[string]float64     `json:"confidence_scores"`
	BoundingBoxes    map[string]interface{} `json:"bounding_boxes"`
	FaceEmbedding    []byte                 `json:"face_embedding,omitempty"`
	ProcessingTimeMS int                    `json:"processing_time_ms"`
	ModelVersion     string                 `json:"model_version"`
}

func (h *DetectionHandler) ProcessFrame(w http.ResponseWriter, r *http.Request) {
	var dto ProcessFrameDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cameraID, err := uuid.Parse(dto.CameraID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid camera_id")
		return
	}

	req := &service.ProcessFrameRequest{
		CameraID:         cameraID,
		FrameURL:         dto.FrameURL,
		DetectedPPE:      dto.DetectedPPE,
		ConfidenceScores: dto.ConfidenceScores,
		BoundingBoxes:    dto.BoundingBoxes,
		FaceEmbedding:    dto.FaceEmbedding,
		ProcessingTimeMS: dto.ProcessingTimeMS,
		ModelVersion:     dto.ModelVersion,
	}

	resp, err := h.detectionSvc.ProcessFrame(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *DetectionHandler) GetDetectionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	detection, err := h.detectionSvc.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "detection not found")
		return
	}

	respondJSON(w, http.StatusOK, detection)
}

func (h *DetectionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	cameraIDStr := r.URL.Query().Get("camera_id")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	cameraID, err := uuid.Parse(cameraIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid camera_id")
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

	stats, err := h.detectionSvc.GetDetectionStats(r.Context(), cameraID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func respondError(w http.ResponseWriter, i int, s string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(i)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": s})
}

func respondJSON(w http.ResponseWriter, i int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(i)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}
