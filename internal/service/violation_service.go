package service

import (
	"context"
	"fmt"
	"time"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/google/uuid"
)

type ViolationService struct {
	violationRepo repository.ViolationRepository
	incidentRepo  repository.IncidentRepository
	userRepo      repository.UserRepository
	camundaClient *CamundaClient
}

func NewViolationService(
	violationRepo repository.ViolationRepository,
	incidentRepo repository.IncidentRepository,
	userRepo repository.UserRepository,
	camundaClient *CamundaClient,
) *ViolationService {
	return &ViolationService{
		violationRepo: violationRepo,
		incidentRepo:  incidentRepo,
		userRepo:      userRepo,
		camundaClient: camundaClient,
	}
}

func (s *ViolationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Violation, error) {
	return s.violationRepo.GetByID(ctx, id)
}

func (s *ViolationService) CreateIncidentFromViolation(ctx context.Context, violationID uuid.UUID) (*models.Incident, error) {
	violation, err := s.violationRepo.GetByID(ctx, violationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get violation: %w", err)
	}

	existingIncident, _ := s.incidentRepo.GetByViolationID(ctx, violationID)
	if existingIncident != nil {
		return existingIncident, nil
	}

	incidentNumber := s.generateIncidentNumber()

	incident := &models.Incident{
		ViolationID:    violationID,
		IncidentNumber: incidentNumber,
		Title:          fmt.Sprintf("PPE Violation: %s", violation.ViolationType),
		Description:    violation.Description,
		Status:         models.IncidentOpen,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if violation.Severity >= models.SeverityHigh {
		supervisor, err := s.userRepo.GetAvailableSupervisor(ctx)
		if err == nil && supervisor != nil {
			incident.AssignedTo = &supervisor.ID
		}
	}

	if err := s.incidentRepo.Create(ctx, incident); err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	processInstanceID, err := s.camundaClient.StartIncidentProcess(ctx, incident)
	if err == nil {
		incident.CamundaProcessInstanceID = &processInstanceID
		s.incidentRepo.Update(ctx, incident)
	}

	s.createIncidentEvent(ctx, incident.ID, "created", nil, "Incident created from violation")

	return incident, nil
}

func (s *ViolationService) AcknowledgeViolation(ctx context.Context, violationID uuid.UUID, acknowledgedBy uuid.UUID) error {
	violation, err := s.violationRepo.GetByID(ctx, violationID)
	if err != nil {
		return err
	}

	now := time.Now()
	violation.AcknowledgedBy = &acknowledgedBy
	violation.AcknowledgedAt = &now

	return s.violationRepo.Update(ctx, violation)
}

func (s *ViolationService) GetViolationStats(ctx context.Context, from, to time.Time) (*ViolationStats, error) {
	violations, err := s.violationRepo.GetByTimeRange(ctx, from, to)
	if err != nil {
		return nil, err
	}

	stats := &ViolationStats{
		Total:      len(violations),
		BySeverity: make(map[models.ViolationSeverity]int),
		ByType:     make(map[string]int),
	}

	for _, v := range violations {
		stats.BySeverity[v.Severity]++
		stats.ByType[v.ViolationType]++
	}

	return stats, nil
}

type ViolationStats struct {
	Total      int                              `json:"total"`
	BySeverity map[models.ViolationSeverity]int `json:"by_severity"`
	ByType     map[string]int                   `json:"by_type"`
}

func (s *ViolationService) generateIncidentNumber() string {
	return fmt.Sprintf("INC-%s", time.Now().Format("20060102-150405"))
}

func (s *ViolationService) createIncidentEvent(ctx context.Context, incidentID uuid.UUID, eventType string, userID *uuid.UUID, comment string) {
	event := &models.IncidentEvent{
		IncidentID: incidentID,
		EventType:  eventType,
		UserID:     userID,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
	_ = s.incidentRepo.CreateEvent(ctx, event)
}
