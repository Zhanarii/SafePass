package service

import (
	"context"
	"fmt"
	"time"

	"ppe-detection/internal/models"
	"ppe-detection/internal/repository"

	"github.com/google/uuid"
)

type ReportService struct {
	violationRepo repository.ViolationRepository
	detectionRepo repository.DetectionRepository
	accessLogRepo repository.AccessLogRepository
	pdfGenerator  *PDFGenerator
}

func NewReportService(
	violationRepo repository.ViolationRepository,
	detectionRepo repository.DetectionRepository,
	accessLogRepo repository.AccessLogRepository,
	pdfGenerator *PDFGenerator,
) *ReportService {
	return &ReportService{
		violationRepo: violationRepo,
		detectionRepo: detectionRepo,
		accessLogRepo: accessLogRepo,
		pdfGenerator:  pdfGenerator,
	}
}

type ViolationReport struct {
	Period          string
	TotalViolations int
	BySeverity      map[models.ViolationSeverity]int
	ByType          map[string]int
	TopViolators    []UserViolationStat
	ComplianceRate  float64
	Violations      []models.Violation
}

type UserViolationStat struct {
	UserID         uuid.UUID
	EmployeeID     string
	FullName       string
	ViolationCount int
}

func (s *ReportService) GenerateViolationReport(ctx context.Context, from, to time.Time) (*ViolationReport, error) {
	violations, err := s.violationRepo.GetByTimeRange(ctx, from, to)
	if err != nil {
		return nil, err
	}

	detections, err := s.detectionRepo.GetByCameraAndTimeRange(ctx, uuid.Nil, from, to, nil)
	if err != nil {
		return nil, err
	}

	report := &ViolationReport{
		Period:          fmt.Sprintf("%s - %s", from.Format("2006-01-02"), to.Format("2006-01-02")),
		TotalViolations: len(violations),
		BySeverity:      make(map[models.ViolationSeverity]int),
		ByType:          make(map[string]int),
		Violations:      violations,
	}

	for _, v := range violations {
		report.BySeverity[v.Severity]++
		report.ByType[v.ViolationType]++
	}

	if len(detections) > 0 {
		compliantCount := 0
		for _, d := range detections {
			if d.Status == models.StatusCompliant {
				compliantCount++
			}
		}
		report.ComplianceRate = float64(compliantCount) / float64(len(detections)) * 100
	}

	return report, nil
}

func (s *ReportService) GeneratePDFReport(ctx context.Context, report *ViolationReport) ([]byte, error) {
	return s.pdfGenerator.Generate(report)
}

type PDFGenerator struct{}

func (g *PDFGenerator) Generate(report *ViolationReport) ([]byte, error) {
	return []byte("PDF content"), nil
}
