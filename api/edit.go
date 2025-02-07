package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
)

type StartEditParams struct {
	ProfileKey      int
	PhotographyType photographyType
	Tools           []Tool
	CallbackURL     url.URL
}

type editStatusResponse struct {
	Data struct {
		Status editStatus `json:"status"`
	} `json:"data"`
}

func (t *client) StartEdit(project Project, params StartEditParams) error {
	if slices.Contains(params.Tools, ToolCrop) && slices.Contains(params.Tools, ToolPortraitCrop) {
		return errors.New("cannot use both crop and portrait crop tools")
	}
	if slices.Contains(params.Tools, ToolPerspectiveCorrection) && slices.Contains(params.Tools, ToolStraighten) {
		return errors.New("cannot use both perspective correction and straighten tools")
	}
	if params.PhotographyType == PhotographyTypeRealEstate && slices.Contains(params.Tools, ToolStraighten) {
		// TODO: warning, do not use straighten tool with real estate photography
		// https://support.imagen-ai.com/hc/en-us/articles/19137253415965-Automate-your-post-production-workflow-with-Imagen-API#h_01J1M5Y4QVAYW9M9422QY64977
	}
	if params.ProfileKey == 0 {
		return errors.New("profile key is required")
	}
	reqData := map[string]any{
		"profile_key": params.ProfileKey,
	}
	if params.CallbackURL.String() != "" {
		reqData["callback_url"] = params.CallbackURL.String()
	}
	if params.PhotographyType != PhotographyTypeNone {
		reqData["photography_type"] = params.PhotographyType
	}
	for _, tool := range params.Tools {
		reqData[string(tool)] = true
	}

	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(reqData); err != nil {
		return fmt.Errorf("failed to encode start edit request: %w", err)
	}

	req, err := t.newPostRequest(fmt.Sprintf("/projects/%s/edit", project.ID), &encoded)
	if err != nil {
		return fmt.Errorf("failed to create start edit request: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start edit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start edit, expected HTTP 200: %s", resp.Status)
	}

	return nil
}

func (t *client) GetEditStatus(project Project) (editStatus, error) {
	req, err := t.newGetRequest(fmt.Sprintf("/projects/%s/edit/status", project.ID))
	if err != nil {
		return EditStatusUnknown, fmt.Errorf("failed to create request to get edit status: %w", err)
	}

	res, err := t.http.Do(req)
	if err != nil {
		return EditStatusUnknown, fmt.Errorf("failed to get edit status: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return EditStatusUnknown, fmt.Errorf("failed to get edit status, expected HTTP 200: %s", res.Status)
	}

	var envelope editStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		return EditStatusUnknown, fmt.Errorf("failed to decode edit status: %w", err)
	}

	return envelope.Data.Status, nil
}
