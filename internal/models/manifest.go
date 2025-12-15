package models

import "time"

type QuestionManifest struct {
	SchemaVersion string            `json:"schema_version"`
	Dataset       string            `json:"dataset"`
	Type          string            `json:"type"`
	Version       string            `json:"version"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Includes      []string          `json:"includes"`
	Counts        map[string]int    `json:"counts"`
	Checksums     map[string]string `json:"checksums"`
}

func NewQuestionManifest(datasetName string) *QuestionManifest {
	return &QuestionManifest{
		SchemaVersion: "qcm/1.0.0",
		Dataset:       datasetName,
		Type:          "questions",
		Version:       "1.0.0",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		Includes: []string{
			"questions",
			"themes",
			"subthemes",
			"tags",
		},
		Counts:    make(map[string]int),
		Checksums: make(map[string]string),
	}
}

type Manifest struct {
	SchemaVersion string            `json:"schema_version"`
	Dataset       string            `json:"dataset"`
	Type          string            `json:"type"`
	Version       string            `json:"version"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Includes      []string          `json:"includes"`
	Counts        map[string]int    `json:"counts"`
	Checksums     map[string]string `json:"checksums"`
}
