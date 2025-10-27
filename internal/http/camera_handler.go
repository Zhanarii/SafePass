package http

import (
	"encoding/json"
	"net/http"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CameraHandler struct {
	cameraRepo repository.CameraRepository
}

func NewCameraHandler(repo repository.CameraRepository) *CameraHandler {
	return &CameraHandler{cameraRepo: repo}
}

type CreateCameraDTO struct {
	LocationID   string `json:"location_id"`
	Name         string `json:"name"`
	RTSPURL      string `json:"rtsp_url"`
	ViewingAngle *int   `json:"viewing_angle,omitempty"`
	FPS          int    `json:"fps"`
	Resolution   string `json:"resolution"`
}

func (h *CameraHandler) CreateCamera(w http.ResponseWriter, r *http.Request) {
	var dto CreateCameraDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	locationID, err := uuid.Parse(dto.LocationID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid location_id")
		return
	}

	camera := &models.Camera{
		LocationID:   locationID,
		Name:         dto.Name,
		RTSPURL:      dto.RTSPURL,
		ViewingAngle: dto.ViewingAngle,
		FPS:          dto.FPS,
		Resolution:   dto.Resolution,
		IsActive:     true,
	}

	if err := h.cameraRepo.Create(r.Context(), camera); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, camera)
}

func (h *CameraHandler) GetCameraByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	camera, err := h.cameraRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "camera not found")
		return
	}

	respondJSON(w, http.StatusOK, camera)
}

func (h *CameraHandler) GetActiveCameras(w http.ResponseWriter, r *http.Request) {
	cameras, err := h.cameraRepo.GetActive(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, cameras)
}

func (h *CameraHandler) UpdateCamera(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	camera, err := h.cameraRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "camera not found")
		return
	}

	var dto CreateCameraDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	camera.Name = dto.Name
	camera.RTSPURL = dto.RTSPURL
	camera.ViewingAngle = dto.ViewingAngle
	camera.FPS = dto.FPS
	camera.Resolution = dto.Resolution

	if err := h.cameraRepo.Update(r.Context(), camera); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, camera)
}
