package service

import (
	"context"
	"fmt"

	"ppe-detection/internal/models"

	"github.com/google/uuid"
)

type NotificationService struct {
	emailClient   *EmailClient
	smsClient     *SMSClient
	webhookClient *WebhookClient
}

func NewNotificationService(
	emailClient *EmailClient,
	smsClient *SMSClient,
	webhookClient *WebhookClient,
) *NotificationService {
	return &NotificationService{
		emailClient:   emailClient,
		smsClient:     smsClient,
		webhookClient: webhookClient,
	}
}

func (s *NotificationService) NotifyViolation(ctx context.Context, violation *models.Violation) error {
	message := fmt.Sprintf(
		"PPE Violation Detected\nType: %s\nSeverity: %s\nTime: %s\nSnapshot: %s",
		violation.ViolationType,
		violation.Severity,
		violation.CreatedAt.Format("2006-01-02 15:04:05"),
		violation.SnapshotURL,
	)

	if violation.Severity >= models.SeverityHigh {
		go s.emailClient.SendAlert("safety@company.com", "Critical PPE Violation", message)
		go s.smsClient.SendAlert("+1234567890", message)
	} else {
		go s.emailClient.SendAlert("safety@company.com", "PPE Violation", message)
	}

	go s.webhookClient.SendViolationAlert(ctx, violation)

	return nil
}

type EmailClient struct{}

func (c *EmailClient) SendAlert(to, subject, body string) error { return nil }

type SMSClient struct{}

func (c *SMSClient) SendAlert(phone, message string) error { return nil }

type WebhookClient struct{}

func (c *WebhookClient) SendViolationAlert(ctx context.Context, v *models.Violation) error {
	return nil
}

type SKUDClient struct{}

func (c *SKUDClient) OpenGate(ctx context.Context, cameraID uuid.UUID, zoneName string) error {
	return nil
}

type CamundaClient struct{}

func (c *CamundaClient) StartIncidentProcess(ctx context.Context, incident *models.Incident) (string, error) {
	return "process-instance-123", nil
}
