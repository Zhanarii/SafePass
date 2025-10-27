// internal/service/detection_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/google/uuid"
)

type DetectionService struct {
	detectionRepo   repository.DetectionRepository
	violationRepo   repository.ViolationRepository
	accessZoneRepo  repository.AccessZoneRepository
	cameraRepo      repository.CameraRepository
	notificationSvc *NotificationService
}

func NewDetectionService(
	detectionRepo repository.DetectionRepository,
	violationRepo repository.ViolationRepository,
	accessZoneRepo repository.AccessZoneRepository,
	cameraRepo repository.CameraRepository,
	notificationSvc *NotificationService,
) *DetectionService {
	return &DetectionService{
		detectionRepo:   detectionRepo,
		violationRepo:   violationRepo,
		accessZoneRepo:  accessZoneRepo,
		cameraRepo:      cameraRepo,
		notificationSvc: notificationSvc,
	}
}

type ProcessFrameRequest struct {
	CameraID         uuid.UUID
	FrameURL         string
	DetectedPPE      []string
	ConfidenceScores map[string]float64
	BoundingBoxes    map[string]interface{}
	FaceEmbedding    []byte
	ProcessingTimeMS int
	ModelVersion     string
}

type ProcessFrameResponse struct {
	DetectionID    uuid.UUID
	Status         models.DetectionStatus
	AccessDecision models.AccessDecision
	Violation      *models.Violation
	Message        string
}

func (s *DetectionService) ProcessFrame(ctx context.Context, req *ProcessFrameRequest) (*ProcessFrameResponse, error) {
	camera, err := s.cameraRepo.GetByID(ctx, req.CameraID)
	if err != nil {
		return nil, fmt.Errorf("failed to get camera: %w", err)
	}

	var accessZone *models.AccessZone
	if camera.AccessZoneID != nil {
		accessZone, err = s.accessZoneRepo.GetByID(ctx, *camera.AccessZoneID)
		if err != nil {
			return nil, fmt.Errorf("failed to get access zone: %w", err)
		}
	}

	missingPPE := s.determineMissingPPE(req.DetectedPPE, accessZone)
	status := s.determineDetectionStatus(missingPPE, req.ConfidenceScores)

	detection := &models.Detection{
		CameraID:         req.CameraID,
		AccessZoneID:     camera.AccessZoneID,
		Timestamp:        time.Now(),
		FrameURL:         req.FrameURL,
		DetectedPPE:      req.DetectedPPE,
		MissingPPE:       missingPPE,
		ConfidenceScores: req.ConfidenceScores,
		BoundingBoxes:    req.BoundingBoxes,
		Status:           status,
		FaceEmbedding:    req.FaceEmbedding,
		ProcessingTimeMS: req.ProcessingTimeMS,
		ModelVersion:     req.ModelVersion,
	}

	if err := s.detectionRepo.Create(ctx, detection); err != nil {
		return nil, fmt.Errorf("failed to create detection: %w", err)
	}

	response := &ProcessFrameResponse{
		DetectionID:    detection.ID,
		Status:         status,
		AccessDecision: models.AccessAllowed,
	}

	if status == models.StatusViolation {
		violation, accessDecision := s.handleViolation(ctx, detection, accessZone)
		response.Violation = violation
		response.AccessDecision = accessDecision
		response.Message = s.buildViolationMessage(missingPPE)

		if violation.Severity == models.SeverityCritical || violation.Severity == models.SeverityHigh {
			go s.notificationSvc.NotifyViolation(context.Background(), violation)
		}
	}

	return response, nil
}

func (s *DetectionService) determineMissingPPE(detected []string, zone *models.AccessZone) []string {
	if zone == nil {
		return []string{}
	}

	detectedMap := make(map[string]bool)
	for _, ppe := range detected {
		detectedMap[ppe] = true
	}

	var missing []string
	for _, required := range zone.RequiredPPE {
		if !detectedMap[required] {
			missing = append(missing, required)
		}
	}

	return missing
}

func (s *DetectionService) determineDetectionStatus(missingPPE []string, confidence map[string]float64) models.DetectionStatus {
	if len(missingPPE) == 0 {
		return models.StatusCompliant
	}

	lowConfidenceCount := 0
	for _, conf := range confidence {
		if conf < 0.7 {
			lowConfidenceCount++
		}
	}

	if lowConfidenceCount > 0 && len(missingPPE) <= 1 {
		return models.StatusWarning
	}

	return models.StatusViolation
}

func (s *DetectionService) handleViolation(
	ctx context.Context,
	detection *models.Detection,
	zone *models.AccessZone,
) (*models.Violation, models.AccessDecision) {
	severity := s.calculateViolationSeverity(detection.MissingPPE, zone)
	accessDecision := s.determineAccessDecision(severity, zone)

	violation := &models.Violation{
		DetectionID:    detection.ID,
		CameraID:       detection.CameraID,
		AccessZoneID:   detection.AccessZoneID,
		ViolationType:  s.buildViolationType(detection.MissingPPE),
		Severity:       severity,
		Description:    s.buildViolationMessage(detection.MissingPPE),
		SnapshotURL:    detection.FrameURL,
		AccessDecision: accessDecision,
		CreatedAt:      time.Now(),
	}

	if err := s.violationRepo.Create(ctx, violation); err != nil {
		fmt.Printf("failed to create violation: %v\n", err)
	}

	return violation, accessDecision
}

func (s *DetectionService) calculateViolationSeverity(missingPPE []string, zone *models.AccessZone) models.ViolationSeverity {
	if zone == nil {
		return models.SeverityLow
	}

	criticalPPE := map[string]bool{
		"helmet": true,
		"mask":   true,
	}

	hasCriticalMissing := false
	for _, ppe := range missingPPE {
		if criticalPPE[ppe] {
			hasCriticalMissing = true
			break
		}
	}

	if hasCriticalMissing && zone.DangerLevel == models.SeverityCritical {
		return models.SeverityCritical
	}

	if len(missingPPE) >= 2 {
		return models.SeverityHigh
	}

	if hasCriticalMissing {
		return models.SeverityHigh
	}

	return models.SeverityMedium
}

func (s *DetectionService) determineAccessDecision(severity models.ViolationSeverity, zone *models.AccessZone) models.AccessDecision {
	if severity == models.SeverityCritical {
		return models.AccessDenied
	}

	if severity == models.SeverityHigh && zone != nil && zone.DangerLevel >= models.SeverityHigh {
		return models.AccessDenied
	}

	if severity == models.SeverityMedium {
		return models.AccessManualReview
	}

	return models.AccessAllowed
}

func (s *DetectionService) buildViolationType(missingPPE []string) string {
	if len(missingPPE) == 0 {
		return "unknown"
	}
	if len(missingPPE) == 1 {
		return fmt.Sprintf("missing_%s", missingPPE[0])
	}
	return "multiple_missing_ppe"
}

func (s *DetectionService) buildViolationMessage(missingPPE []string) string {
	if len(missingPPE) == 0 {
		return "No violations detected"
	}
	return fmt.Sprintf("Missing required PPE: %v", missingPPE)
}

func (s *DetectionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Detection, error) {
	return s.detectionRepo.GetByID(ctx, id)
}

func (s *DetectionService) GetDetectionStats(ctx context.Context, cameraID uuid.UUID, from, to time.Time) (*DetectionStats, error) {
	detections, err := s.detectionRepo.GetByCameraAndTimeRange(ctx, cameraID, from, to, nil)
	if err != nil {
		return nil, err
	}

	stats := &DetectionStats{
		TotalDetections: len(detections),
		CompliantCount:  0,
		ViolationCount:  0,
		WarningCount:    0,
	}

	for _, d := range detections {
		switch d.Status {
		case models.StatusCompliant:
			stats.CompliantCount++
		case models.StatusViolation:
			stats.ViolationCount++
		case models.StatusWarning:
			stats.WarningCount++
		}
	}

	if stats.TotalDetections > 0 {
		stats.ComplianceRate = float64(stats.CompliantCount) / float64(stats.TotalDetections) * 100
	}

	return stats, nil
}

type DetectionStats struct {
	TotalDetections int     `json:"total_detections"`
	CompliantCount  int     `json:"compliant_count"`
	ViolationCount  int     `json:"violation_count"`
	WarningCount    int     `json:"warning_count"`
	ComplianceRate  float64 `json:"compliance_rate"`
}
