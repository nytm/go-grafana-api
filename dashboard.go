package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// DashboardMeta holds dashboard metadata
type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
}

// DashboardSaveResponse represents the response from the API when
// a dashboard is saved
type DashboardSaveResponse struct {
	ID      int64  `json:"id"`
	UID     int64  `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
	Slug    string `json:"slug"`
}

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	Meta  DashboardMeta          `json:"meta"`
	Model map[string]interface{} `json:"dashboard"`
}

// SaveDashboard saves the given dashboard model to the API
func (c *Client) SaveDashboard(model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/api/dashboards/db", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// Dashboard gets the dashboard with the given URI from Grafana
func (c *Client) Dashboard(uri string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/dashboards/%s", uri)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Dashboard{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// DeleteDashboard will delete the dashboard with the given slug from Grafana
func (c *Client) DeleteDashboard(slug string) error {
	path := fmt.Sprintf("/api/dashboards/db/%s", slug)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}
