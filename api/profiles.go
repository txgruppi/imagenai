package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type Profile struct {
	Key         int         `json:"profile_key"`
	Name        string      `json:"profile_name"`
	Type        string      `json:"profile_type"`
	ImageFormat imageFormat `json:"image_type"`
}

type Profiles []Profile

func (t Profiles) Len() int {
	return len(t)
}

func (t Profiles) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t Profiles) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type profilesResponseEnvelope struct {
	Data struct {
		Profiles Profiles `json:"profiles"`
	} `json:"data"`
}

func (t *client) GetAvailableProfiles() (Profiles, error) {
	req, err := t.newGetRequest("/profiles")
	if err != nil {
		return nil, fmt.Errorf("failed to create request to get available profiles: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get available profiles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get available profiles, expected HTTP 200: %s", resp.Status)
	}

	var envelope profilesResponseEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("failed to decode available profiles: %w", err)
	}

	sort.Sort(envelope.Data.Profiles)
	return envelope.Data.Profiles, nil
}
