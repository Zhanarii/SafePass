package service

import (
	"context"
	"fmt"
	"time"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/google/uuid"
)

type AccessControlService struct {
	accessLogRepo  repository.AccessLogRepository
	detectionRepo  repository.DetectionRepository
	accessZoneRepo repository.AccessZoneRepository
	userRepo       repository.UserRepository
	skudClient     *SKUDClient
}

func NewAccessControlService(
	accessLogRepo repository.AccessLogRepository,
	detectionRepo repository.DetectionRepository,
	accessZoneRepo repository.AccessZoneRepository,
	userRepo repository.UserRepository,
	skudClient *SKUDClient,
) *AccessControlService {
	return &AccessControlService{
		accessLogRepo:  accessLogRepo,
		detectionRepo:  detectionRepo,
		accessZoneRepo: accessZoneRepo,
		userRepo:       userRepo,
		skudClient:     skudClient,
	}
}

type AccessRequest struct {
	UserID       *uuid.UUID
	BadgeNumber  *string
	CameraID     uuid.UUID
	AccessZoneID uuid.UUID
	DetectionID  *uuid.UUID
}

type AccessResponse struct {
	Decision  models.AccessDecision
	Message   string
	LogID     uuid.UUID
	AllowedAt *time.Time
}

func (s *AccessControlService) CheckAccess(ctx context.Context, req *AccessRequest) (*AccessResponse, error) {
	zone, err := s.accessZoneRepo.GetByID(ctx, req.AccessZoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access zone: %w", err)
	}

	if !zone.IsActive {
		return s.denyAccess(ctx, req, "Access zone is inactive")
	}

	var detection *models.Detection
	if req.DetectionID != nil {
		detection, err = s.detectionRepo.GetByID(ctx, *req.DetectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get detection: %w", err)
		}
	}

	var decision models.AccessDecision
	var message string

	if detection != nil {
		if detection.Status == models.StatusCompliant {
			decision = models.AccessAllowed
			message = "All required PPE detected"
		} else if detection.Status == models.StatusViolation {
			if len(detection.MissingPPE) > 0 {
				decision = models.AccessDenied
				message = fmt.Sprintf("Missing required PPE: %v", detection.MissingPPE)
			} else {
				decision = models.AccessManualReview
				message = "PPE detection uncertain, manual review required"
			}
		} else {
			decision = models.AccessManualReview
			message = "Warning status, manual review required"
		}
	} else {
		decision = models.AccessManualReview
		message = "No PPE detection available"
	}

	now := time.Now()
	accessLog := &models.AccessLog{
		UserID:       req.UserID,
		CameraID:     req.CameraID,
		AccessZoneID: &req.AccessZoneID,
		DetectionID:  req.DetectionID,
		Decision:     decision,
		Timestamp:    now,
		BadgeScanned: req.BadgeNumber != nil,
	}

	if err := s.accessLogRepo.Create(ctx, accessLog); err != nil {
		return nil, fmt.Errorf("failed to create access log: %w", err)
	}

	var allowedAt *time.Time
	if decision == models.AccessAllowed {
		if err := s.skudClient.OpenGate(ctx, req.CameraID, zone.Name); err != nil {
			fmt.Printf("failed to open gate: %v\n", err)
		}
		allowedAt = &now
	}

	return &AccessResponse{
		Decision:  decision,
		Message:   message,
		LogID:     accessLog.ID,
		AllowedAt: allowedAt,
	}, nil
}

func (s *AccessControlService) OverrideAccess(
	ctx context.Context,
	logID uuid.UUID,
	overrideBy uuid.UUID,
	reason string,
) error {
	log, err := s.accessLogRepo.GetByID(ctx, logID)
	if err != nil {
		return fmt.Errorf("failed to get access log: %w", err)
	}

	log.OverrideBy = &overrideBy
	log.OverrideReason = &reason
	log.Decision = models.AccessAllowed

	if err := s.accessLogRepo.Update(ctx, log); err != nil {
		return fmt.Errorf("failed to update access log: %w", err)
	}

	zone, _ := s.accessZoneRepo.GetByID(ctx, *log.AccessZoneID)
	if zone != nil {
		s.skudClient.OpenGate(ctx, log.CameraID, zone.Name)
	}

	return nil
}

func (s *AccessControlService) denyAccess(ctx context.Context, req *AccessRequest, reason string) (*AccessResponse, error) {
	accessLog := &models.AccessLog{
		UserID:       req.UserID,
		CameraID:     req.CameraID,
		AccessZoneID: &req.AccessZoneID,
		DetectionID:  req.DetectionID,
		Decision:     models.AccessDenied,
		Timestamp:    time.Now(),
	}

	if err := s.accessLogRepo.Create(ctx, accessLog); err != nil {
		return nil, err
	}

	return &AccessResponse{
		Decision: models.AccessDenied,
		Message:  reason,
		LogID:    accessLog.ID,
	}, nil
}

func (s *AccessControlService) GetAccessHistory(
	ctx context.Context,
	userID uuid.UUID,
	from, to time.Time,
) ([]models.AccessLog, error) {
	return s.accessLogRepo.GetByUserAndTimeRange(ctx, userID, from, to)
}
