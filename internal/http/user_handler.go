package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

type CreateUserDTO struct {
	EmployeeID  string  `json:"employee_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone,omitempty"`
	Role        string  `json:"role"`
	Department  *string `json:"department,omitempty"`
	BadgeNumber *string `json:"badge_number,omitempty"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := &models.User{
		EmployeeID:  dto.EmployeeID,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Email:       dto.Email,
		Phone:       dto.Phone,
		Role:        models.UserRole(dto.Role),
		Department:  dto.Department,
		BadgeNumber: dto.BadgeNumber,
		IsActive:    true,
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	users, err := h.userRepo.GetAll(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Email = dto.Email
	user.Phone = dto.Phone
	user.Role = models.UserRole(dto.Role)
	user.Department = dto.Department
	user.BadgeNumber = dto.BadgeNumber

	if err := h.userRepo.Update(r.Context(), user); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.userRepo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
