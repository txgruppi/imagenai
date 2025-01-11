package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	ID string `json:"project_uuid"`
}

type projectResponseEnvelope struct {
	Data Project `json:"data"`
}

func (t *client) CreateProject() (Project, error) {
	// The trailing slash is required. I know, it sucks.
	req, err := t.newPostRequest("/projects/", http.NoBody)
	if err != nil {
		return Project{}, fmt.Errorf("failed to create request to create project: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return Project{}, fmt.Errorf("failed to create project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Project{}, fmt.Errorf("failed to create project, expected HTTP 200: %s", resp.Status)
	}

	var envelope projectResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return Project{}, fmt.Errorf("failed to decode project: %w", err)
	}

	return envelope.Data, nil
}
