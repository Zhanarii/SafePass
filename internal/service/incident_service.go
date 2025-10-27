package service

import (
	"context"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/google/uuid"
)

type IncidentService struct {
	incidentRepo repository.IncidentRepository
}

func NewIncidentService(incidentRepo repository.IncidentRepository) *IncidentService {
	return &IncidentService{incidentRepo: incidentRepo}
}

func (s *IncidentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	return s.incidentRepo.GetByID(ctx, id)
}

func (s *IncidentService) GetByStatus(ctx context.Context, status models.IncidentStatus) ([]models.Incident, error) {
	return s.incidentRepo.GetByStatus(ctx, status)
}

func (s *IncidentService) Update(ctx context.Context, incident *models.Incident) error {
	return s.incidentRepo.Update(ctx, incident)
}
