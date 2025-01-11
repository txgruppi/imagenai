package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExportDownload struct {
	FileName     string `json:"file_name"`
	DownloadLink string `json:"download_link"`
}

type exportDownloadEnvelope struct {
	Data struct {
		Files []ExportDownload `json:"files_list"`
	} `json:"data"`
}

type exportStatusResponse struct {
	Data struct {
		Status exportStatus `json:"status"`
	} `json:"data"`
}

func (t *client) StartExport(project Project) error {
	req, err := t.newPostRequest(fmt.Sprintf("/projects/%s/export", project.ID), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create start export request: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start export: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start export, expected HTTP 200: %s", resp.Status)
	}

	return nil
}

func (t *client) GetExportStatus(project Project) (exportStatus, error) {
	req, err := t.newGetRequest(fmt.Sprintf("/projects/%s/export/status", project.ID))
	if err != nil {
		return ExportStatusUnknown, fmt.Errorf("failed to create request to get export status: %w", err)
	}

	res, err := t.http.Do(req)
	if err != nil {
		return ExportStatusUnknown, fmt.Errorf("failed to get export status: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ExportStatusUnknown, fmt.Errorf("failed to get export status, expected HTTP 200: %s", res.Status)
	}

	var envelope exportStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return ExportStatusUnknown, fmt.Errorf("failed to decode export status: %w", err)
	}

	return envelope.Data.Status, nil
}

func (t *client) GetExportDownloadLinks(project Project) ([]ExportDownload, error) {
	req, err := t.newGetRequest(fmt.Sprintf("/projects/%s/export/get_temporary_download_links", project.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to get export download links: %w", err)
	}

	res, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get export download links: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get export download links, expected HTTP 200: %s", res.Status)
	}

	var envelope exportDownloadEnvelope
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode export download links: %w", err)
	}

	return envelope.Data.Files, nil
}
