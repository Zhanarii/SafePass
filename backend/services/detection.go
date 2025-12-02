package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Detector struct {
	ModelPath string
}

type Detection struct {
	Class      string    `json:"class"`
	Confidence float64   `json:"confidence"`
	BBox       []float64 `json:"bbox"`
}

type DetectionResult struct {
	Detections    []Detection `json:"detections"`
	HasViolations bool        `json:"has_violations"`
	Violations    []Detection `json:"violations"`
}

func NewDetector(modelPath string) (*Detector, error) {
	return &Detector{
		ModelPath: modelPath,
	}, nil
}

func (d *Detector) DetectFromFile(imagePath string) (*DetectionResult, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	return d.DetectViaAPI(imageData)
}

func (d *Detector) DetectViaAPI(imageData []byte) (*DetectionResult, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, bytes.NewReader(imageData)); err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://localhost:8001/detect", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ML service not available: %w", err)
	}
	defer resp.Body.Close()

	var result DetectionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
