package report

import (
	"encoding/json"
	"time"
)

type Report struct {
	Timestamp   time.Time   `json:"timestamp"`
	Image       string      `json:"image"`
	ContainerID string      `json:"container_id"`
	Result      interface{} `json:"result"`
}

func BuildReport(image, containerID string, result interface{}) *Report {
	return &Report{
		Timestamp:   time.Now(),
		Image:       image,
		ContainerID: containerID,
		Result:      result,
	}
}

func ToJSON(r *Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
